package meowbail

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// SendStickerPackFromMedia mengunggah media stiker yang dikirim/di-reply dan membuat StickerPackMessage resmi WhatsApp
func (c *Client) SendStickerPackFromMedia(ctx context.Context, chat types.JID, stickerWebpData []byte, packName, publisher string) error {
	if len(stickerWebpData) == 0 {
		return fmt.Errorf("data stiker kosong")
	}

	if packName == "" {
		packName = "Dongtube Pack"
	}
	if publisher == "" {
		publisher = "Dongtube"
	}

	// Upload stiker media ke server WhatsApp
	uploaded, err := c.UploadMedia(ctx, stickerWebpData, whatsmeow.MediaImage)
	if err != nil {
		return fmt.Errorf("upload sticker pack: %w", err)
	}

	h := sha256.Sum256(stickerWebpData)
	fileName := fmt.Sprintf("%x.webp", h[:16])

	origin := waE2E.StickerPackMessage_USER_CREATED
	disappearingModeInitiator := waE2E.DisappearingMode_CHANGED_IN_CHAT

	msg := &waE2E.Message{
		StickerPackMessage: &waE2E.StickerPackMessage{
			Stickers: []*waE2E.StickerPackMessage_Sticker{
				{
					Emojis:             []string{""},
					FileName:           proto.String(fileName),
					IsAnimated:         proto.Bool(false),
					AccessibilityLabel: proto.String(packName),
					IsLottie:           proto.Bool(false),
					Mimetype:           proto.String("image/webp"),
				},
			},
			StickerPackID:       proto.String(fmt.Sprintf("Pack_%x", h[:8])),
			Name:                proto.String(packName),
			Publisher:           proto.String(publisher),
			FileLength:          proto.Uint64(uploaded.FileLength),
			FileSHA256:          uploaded.FileSHA256,
			FileEncSHA256:       uploaded.FileEncSHA256,
			MediaKey:            uploaded.MediaKey,
			DirectPath:          proto.String(uploaded.DirectPath),
			PackDescription:     proto.String("Sticker pack dari " + packName),
			MediaKeyTimestamp:   proto.Int64(time.Now().Unix()),
			TrayIconFileName:    proto.String("tray_icon.webp"),
			ThumbnailDirectPath: proto.String(uploaded.DirectPath),
			ThumbnailSHA256:     uploaded.FileSHA256,
			ThumbnailEncSHA256:  uploaded.FileEncSHA256,
			ThumbnailHeight:     proto.Uint32(252),
			ThumbnailWidth:      proto.Uint32(252),
			ImageDataHash:       proto.String(fmt.Sprintf("%x", h[:])),
			StickerPackSize:     proto.Uint64(uploaded.FileLength),
			StickerPackOrigin:   &origin,
			ContextInfo: &waE2E.ContextInfo{
				IsForwarded:     proto.Bool(true),
				ForwardingScore: proto.Uint32(1),
				Expiration:      proto.Uint32(86400),
				DisappearingMode: &waE2E.DisappearingMode{
					Initiator: &disappearingModeInitiator,
				},
			},
		},
	}

	_, err = c.Client.SendMessage(ctx, chat, msg)
	return err
}
