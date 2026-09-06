package meowbail

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os/exec"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// SendStickerPackMultipleMedia mengemas kumpulan media stiker ke dalam zip format WhatsApp sticker pack resmi dan mengirim StickerPackMessage
func (c *Client) SendStickerPackMultipleMedia(ctx context.Context, chat types.JID, stickerItems [][]byte, packName, publisher string) error {
	if len(stickerItems) == 0 {
		return fmt.Errorf("daftar stiker kosong")
	}

	if packName == "" {
		packName = "Dongtube Pack"
	}
	if publisher == "" {
		publisher = "Dongtube"
	}

	packHash := sha256.Sum256([]byte(fmt.Sprintf("%s_%d", packName, time.Now().UnixNano())))
	packID := fmt.Sprintf("Pack_%x", packHash[:8])
	trayIconFileName := "tray_icon.webp"

	// Build ZIP container (metode Store / uncompressed) sesuai spesifikasi WA sticker pack
	zipBuf := new(bytes.Buffer)
	zw := zip.NewWriter(zipBuf)

	var stickers []*waE2E.StickerPackMessage_Sticker

	for i, itemData := range stickerItems {
		h := sha256.Sum256(itemData)
		// Gunakan RawURLEncoding (tanpa padding '=') persis seperti WhatsApp client
		b64Hash := base64.RawURLEncoding.EncodeToString(h[:])
		fileName := b64Hash + ".webp"

		w, err := zw.CreateHeader(&zip.FileHeader{
			Name:   fileName,
			Method: zip.Store,
		})
		if err != nil {
			return fmt.Errorf("create zip entry %s: %w", fileName, err)
		}
		if _, err := w.Write(itemData); err != nil {
			return fmt.Errorf("write zip entry %s: %w", fileName, err)
		}

		stickers = append(stickers, &waE2E.StickerPackMessage_Sticker{
			Emojis:             []string{""},
			FileName:           proto.String(fileName),
			IsAnimated:         proto.Bool(false),
			AccessibilityLabel: proto.String(fmt.Sprintf("%s #%d", packName, i+1)),
			IsLottie:           proto.Bool(false),
			Mimetype:           proto.String("image/webp"),
		})
	}

	// Buat tray_icon.webp 96x96 sesuai ukuran standar WhatsApp tray icon
	cmdTray := exec.Command("ffmpeg", "-i", "pipe:0", "-vf", "scale=96:96:force_original_aspect_ratio=decrease,pad=96:96:(96-iw)/2:(96-ih)/2:color=0x00000000", "-vcodec", "libwebp", "-f", "webp", "pipe:1")
	cmdTray.Stdin = bytes.NewReader(stickerItems[0])
	trayBytes, err := cmdTray.Output()
	if err != nil || len(trayBytes) == 0 {
		trayBytes = stickerItems[0]
	}

	// Masukkan tray cover ke zip dengan nama baku WhatsApp: tray_icon.webp
	wTray, err := zw.CreateHeader(&zip.FileHeader{
		Name:   trayIconFileName,
		Method: zip.Store,
	})
	if err != nil {
		return fmt.Errorf("create zip tray: %w", err)
	}
	if _, err := wTray.Write(trayBytes); err != nil {
		return fmt.Errorf("write zip tray: %w", err)
	}

	if err := zw.Close(); err != nil {
		return fmt.Errorf("close zip writer: %w", err)
	}

	zipBytes := zipBuf.Bytes()

	// 1. Upload ZIP stiker pack dengan MediaType resmi MediaStickerPack
	uploadedZip, err := c.UploadMedia(ctx, zipBytes, whatsmeow.MediaStickerPack)
	if err != nil {
		uploadedZip, err = c.UploadMedia(ctx, zipBytes, whatsmeow.MediaDocument)
		if err != nil {
			return fmt.Errorf("upload zip sticker pack: %w", err)
		}
	}

	// Generate 252x252 JPEG thumbnail untuk tray icon
	cmdThumb := exec.Command("ffmpeg", "-i", "pipe:0", "-vf", "scale=252:252:force_original_aspect_ratio=decrease,pad=252:252:(252-iw)/2:(252-ih)/2:color=0x00000000", "-vcodec", "mjpeg", "-f", "image2", "pipe:1")
	cmdThumb.Stdin = bytes.NewReader(stickerItems[0])
	thumbJpeg, err := cmdThumb.Output()
	if err != nil || len(thumbJpeg) == 0 {
		cmdThumbFallback := exec.Command("ffmpeg", "-i", "pipe:0", "-vf", "scale=252:252", "-vcodec", "mjpeg", "-f", "image2", "pipe:1")
		cmdThumbFallback.Stdin = bytes.NewReader(stickerItems[0])
		thumbJpeg, _ = cmdThumbFallback.Output()
	}

	var thumbDirectPath string
	var thumbSha256 []byte
	var thumbEncSha256 []byte
	var imageDataHash string

	// 2. Upload thumbnail JPEG (MediaLinkThumbnail)
	if len(thumbJpeg) > 0 {
		hThumb := sha256.Sum256(thumbJpeg)
		imageDataHash = base64.StdEncoding.EncodeToString(hThumb[:])

		uploadedThumb, errThumb := c.UploadMedia(ctx, thumbJpeg, whatsmeow.MediaLinkThumbnail)
		if errThumb == nil && uploadedThumb != nil {
			thumbDirectPath = uploadedThumb.DirectPath
			thumbSha256 = uploadedThumb.FileSHA256
			thumbEncSha256 = uploadedThumb.FileEncSHA256
		}
	}

	if thumbDirectPath == "" {
		thumbDirectPath = uploadedZip.DirectPath
		thumbSha256 = uploadedZip.FileSHA256
		thumbEncSha256 = uploadedZip.FileEncSHA256
		if imageDataHash == "" {
			hZip := sha256.Sum256(zipBytes)
			imageDataHash = base64.StdEncoding.EncodeToString(hZip[:])
		}
	}

	origin := waE2E.StickerPackMessage_USER_CREATED
	disappearingModeInitiator := waE2E.DisappearingMode_CHANGED_IN_CHAT

	msg := &waE2E.Message{
		StickerPackMessage: &waE2E.StickerPackMessage{
			Stickers:            stickers,
			StickerPackID:       proto.String(packID),
			Name:                proto.String(packName),
			Publisher:           proto.String(publisher),
			FileLength:          proto.Uint64(uploadedZip.FileLength),
			FileSHA256:          uploadedZip.FileSHA256,
			FileEncSHA256:       uploadedZip.FileEncSHA256,
			MediaKey:            uploadedZip.MediaKey,
			DirectPath:          proto.String(uploadedZip.DirectPath),
			PackDescription:     proto.String(fmt.Sprintf("%s (%d stickers)", packName, len(stickers))),
			MediaKeyTimestamp:   proto.Int64(time.Now().Unix()),
			TrayIconFileName:    proto.String(trayIconFileName),
			ThumbnailDirectPath: proto.String(thumbDirectPath),
			ThumbnailSHA256:     thumbSha256,
			ThumbnailEncSHA256:  thumbEncSha256,
			ThumbnailHeight:     proto.Uint32(252),
			ThumbnailWidth:      proto.Uint32(252),
			ImageDataHash:       proto.String(imageDataHash),
			StickerPackSize:     proto.Uint64(uint64(len(zipBytes))),
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

// SendStickerPackFromMedia wrapper backward-compatibility untuk 1 stiker
func (c *Client) SendStickerPackFromMedia(ctx context.Context, chat types.JID, stickerWebpData []byte, packName, publisher string) error {
	return c.SendStickerPackMultipleMedia(ctx, chat, [][]byte{stickerWebpData}, packName, publisher)
}
