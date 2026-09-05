package meowbail

import (
	"context"
	"encoding/base64"

	"go.mau.fi/whatsmeow"
	waBinary "go.mau.fi/whatsmeow/binary"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// SendExactDongtubeStickerPack mengirim exact Sticker Pack Message dari payload spesifik WhatsApp
func (c *Client) SendExactDongtubeStickerPack(ctx context.Context, chat types.JID) error {
	fileSha256, _ := base64.StdEncoding.DecodeString("kaBpHRYRRQ7ZMYM/nWOjCjnKCsiQS7IXteFhtd6+I6Y=")
	fileEncSha256, _ := base64.StdEncoding.DecodeString("mthbdUKbsu4BBfVP3CI/qQ/xPV0Pva96XnTtCnmCrFE=")
	mediaKey, _ := base64.StdEncoding.DecodeString("0xhKvU7MjLV4e1V3L0fL1A5bIZSvL2RdNeeBI2hn8Po=")
	thumbnailSha256, _ := base64.StdEncoding.DecodeString("UB604LEos5y0UDjeIhREmk+6uky/pGFAbEGGhM0FdvA=")
	thumbnailEncSha256, _ := base64.StdEncoding.DecodeString("/EWG65S7ViLEkioHVUI8n+uSrdgqBb75u8YPtFoKeXU=")

	origin := waE2E.StickerPackMessage_USER_CREATED
	disappearingModeInitiator := waE2E.DisappearingMode_CHANGED_IN_CHAT

	msg := &waE2E.Message{
		StickerPackMessage: &waE2E.StickerPackMessage{
			Stickers: []*waE2E.StickerPackMessage_Sticker{
				{
					Emojis:             []string{""},
					FileName:           proto.String("KDBrwiFN1f6DqiZXqav4qIKRIjm5k-m3sUx0LVqzSsc.webp"),
					IsAnimated:         proto.Bool(false),
					AccessibilityLabel: proto.String(""),
					IsLottie:           proto.Bool(false),
					Mimetype:           proto.String("image/webp"),
				},
				{
					Emojis:             []string{""},
					FileName:           proto.String("HwPjo2O718Xryfvg7zLiTYbBrxza_v0nCf0ud54F-Kw.webp"),
					IsAnimated:         proto.Bool(false),
					AccessibilityLabel: proto.String(""),
					IsLottie:           proto.Bool(false),
					Mimetype:           proto.String("image/webp"),
				},
			},
			StickerPackID:       proto.String("Pack_ba978da79e5ca58b"),
			Name:                proto.String("Dongtube"),
			Publisher:           proto.String("Al"),
			FileLength:          proto.Uint64(54758),
			FileSHA256:          fileSha256,
			FileEncSHA256:       fileEncSha256,
			MediaKey:            mediaKey,
			DirectPath:          proto.String("/v/t62.15575-24/797445000_1099705659182285_5237502063562763067_n.enc?ccb=11-4&oh=01_Q5Aa5gEnqsvPgjv_5eiWokXzaENh4scDwI3tmyDKrP18hmxvpA&oe=6AC36FD9&_nc_sid=5e03e0"),
			PackDescription:     proto.String("Sticker pack dari album"),
			MediaKeyTimestamp:   proto.Int64(1788613023),
			TrayIconFileName:    proto.String("tray_icon.webp"),
			ThumbnailDirectPath: proto.String("/v/t62.15575-24/795553422_2669799446748590_6586724882691871289_n.enc?ccb=11-4&oh=01_Q5Aa5gEtA47MgE83C_3XG8gjRd33djKlHy4Nd4lgYEzafmvHmA&oe=6AC397AF&_nc_sid=5e03e0"),
			ThumbnailSHA256:     thumbnailSha256,
			ThumbnailEncSHA256:  thumbnailEncSha256,
			ThumbnailHeight:     proto.Uint32(252),
			ThumbnailWidth:      proto.Uint32(252),
			ImageDataHash:       proto.String("UB604LEos5y0UDjeIhREmk+6uky/pGFAbEGGhM0FdvA="),
			StickerPackSize:     proto.Uint64(54758),
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

	attrs := waBinary.Attrs{"type": "media"}
	_, err := c.Client.SendMessage(ctx, chat, msg, whatsmeow.SendRequestExtra{
		AdditionalNodes: &[]waBinary.Node{{
			Tag:   "additional_attributes",
			Attrs: attrs,
		}},
	})
	return err
}
