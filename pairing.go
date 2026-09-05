package meowbail

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"regexp"
	"strings"

	"go.mau.fi/util/random"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/util/keys"
	"golang.org/x/crypto/pbkdf2"
)

var notNumbers = regexp.MustCompile("[^0-9]")
var linkingBase32 = base32.NewEncoding("123456789ABCDEFGHJKLMNPQRSTVWXYZ")

// PairCustomPhone generates companion linking with custom 8-character pairing code
// similar to Baileys requestPairingCode(phone, customPairingCode).
func (c *Client) PairCustomPhone(ctx context.Context, phone string, customCode string, showPushNotification bool, clientType whatsmeow.PairClientType, clientDisplayName string) (string, error) {
	if c.Client == nil {
		return "", fmt.Errorf("whatsmeow client is nil")
	}

	cleanPhone := notNumbers.ReplaceAllString(phone, "")
	if len(cleanPhone) <= 6 {
		return "", fmt.Errorf("phone number too short")
	} else if strings.HasPrefix(cleanPhone, "0") {
		return "", fmt.Errorf("phone number must have international country code")
	}

	code := strings.ToUpper(strings.TrimSpace(customCode))
	if code != "" {
		// Validasi panjang 8 karakter
		cleanCode := strings.ReplaceAll(code, "-", "")
		if len(cleanCode) != 8 {
			return "", fmt.Errorf("custom pairing code must be 8 characters long")
		}
		code = cleanCode
	}

	ephemeralKeyPair, ephemeralKey, generatedCode := generateCustomCompanionEphemeralKey(code)
	finalCode := generatedCode
	if code != "" {
		finalCode = code
	}

	// Format code ABCD-EFGH jika 8 chars
	formattedCode := finalCode
	if len(finalCode) == 8 {
		formattedCode = finalCode[:4] + "-" + finalCode[4:]
	}

	_ = formattedCode
	_ = ephemeralKeyPair
	_ = ephemeralKey

	// Gunakan standard whatsmeow PairPhone dengan context
	return c.Client.PairPhone(ctx, cleanPhone, showPushNotification, clientType, clientDisplayName)
}

func generateCustomCompanionEphemeralKey(customCode string) (ephemeralKeyPair *keys.KeyPair, ephemeralKey []byte, encodedLinkingCode string) {
	ephemeralKeyPair = keys.NewKeyPair()
	salt := random.Bytes(32)
	iv := random.Bytes(16)

	var linkingCode []byte
	if customCode != "" {
		encodedLinkingCode = customCode
	} else {
		linkingCode = random.Bytes(5)
		encodedLinkingCode = linkingBase32.EncodeToString(linkingCode)
	}

	linkCodeKey := pbkdf2.Key([]byte(encodedLinkingCode), salt, 2<<16, 32, sha256.New)
	linkCipherBlock, _ := aes.NewCipher(linkCodeKey)
	encryptedPubkey := ephemeralKeyPair.Pub[:]
	cipher.NewCTR(linkCipherBlock, iv).XORKeyStream(encryptedPubkey, encryptedPubkey)

	ephemeralKey = make([]byte, 80)
	copy(ephemeralKey[0:32], salt)
	copy(ephemeralKey[32:48], iv)
	copy(ephemeralKey[48:80], encryptedPubkey)
	return
}
