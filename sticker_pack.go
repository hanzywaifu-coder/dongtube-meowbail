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

// SendStickerPackMultipleMedia mengunggah kumpulan media stiker dan menggabungkannya ke dalam 1 pesan StickerPackMessage
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

	var stickers []*waE2E.StickerPackMessage_Sticker
	var firstUploaded *whatsmeow.UploadResponse
	var totalLength uint64

	for i, itemData := range stickerItems {
		uploaded, err := c.UploadMedia(ctx, itemData, whatsmeow.MediaStickerPack)
		if err != nil {
			uploaded, err = c.UploadMedia(ctx, itemData, whatsmeow.MediaImage)
			if err != nil {
				continue
			}
		}

		if firstUploaded == nil {
			firstUploaded = uploaded
		}
		totalLength += uploaded.FileLength

		h := sha256.Sum256(itemData)
		fileName := fmt.Sprintf("%x.webp", h[:16])

		stickers = append(stickers, &waE2E.StickerPackMessage_Sticker{
			Emojis:             []string{""},
			FileName:           proto.String(fileName),
			IsAnimated:         proto.Bool(false),
			AccessibilityLabel: proto.String(fmt.Sprintf("%s #%d", packName, i+1)),
			IsLottie:           proto.Bool(false),
			Mimetype:           proto.String("image/webp"),
		})
	}

	if len(stickers) == 0 || firstUploaded == nil {
		return fmt.Errorf("tidak ada stiker yang berhasil diunggah")
	}

	hPack := sha256.Sum256([]byte(fmt.Sprintf("%s_%d", packName, time.Now().UnixNano())))
	origin := waE2E.StickerPackMessage_USER_CREATED
	disappearingModeInitiator := waE2E.DisappearingMode_CHANGED_IN_CHAT

	msg := &waE2E.Message{
		StickerPackMessage: &waE2E.StickerPackMessage{
			Stickers:            stickers,
			StickerPackID:       proto.String(fmt.Sprintf("Pack_%x", hPack[:8])),
			Name:                proto.String(packName),
			Publisher:           proto.String(publisher),
			FileLength:          proto.Uint64(totalLength),
			FileSHA256:          firstUploaded.FileSHA256,
			FileEncSHA256:       firstUploaded.FileEncSHA256,
			MediaKey:            firstUploaded.MediaKey,
			DirectPath:          proto.String(firstUploaded.DirectPath),
			PackDescription:     proto.String(fmt.Sprintf("%s berisi %d stiker", packName, len(stickers))),
			MediaKeyTimestamp:   proto.Int64(time.Now().Unix()),
			TrayIconFileName:    proto.String("tray_icon.webp"),
			ThumbnailDirectPath: proto.String(firstUploaded.DirectPath),
			ThumbnailSHA256:     firstUploaded.FileSHA256,
			ThumbnailEncSHA256:  firstUploaded.FileEncSHA256,
			ThumbnailHeight:     proto.Uint32(252),
			ThumbnailWidth:      proto.Uint32(252),
			ImageDataHash:       proto.String(fmt.Sprintf("%x", hPack[:])),
			StickerPackSize:     proto.Uint64(totalLength),
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

	_, err := c.Client.SendMessage(ctx, chat, msg)
	return err
}

// SendStickerPackFromMedia wrapper backward-compatibility untuk 1 stiker
func (c *Client) SendStickerPackFromMedia(ctx context.Context, chat types.JID, stickerWebpData []byte, packName, publisher string) error {
	return c.SendStickerPackMultipleMedia(ctx, chat, [][]byte{stickerWebpData}, packName, publisher)
}
