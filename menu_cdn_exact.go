package meowbail

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// SendRawCDNMenu mengirim pesan persis tanpa ubahan 1 bait pun dari https://cdn.dongtube.id/7e716445aee7.bin
func (c *Client) SendRawCDNMenu(ctx context.Context, chat types.JID, thumbBytes []byte) error {
	secretBytes, _ := base64.StdEncoding.DecodeString("IXM8c2x6FioXJZLibSUFkhheds8R4KQtoWqKWvhwIkY=")

	catParamsJSON := `{"title":"\u0000","sections":[{"title":"ᴋᴀᴛᴇɢᴏʀɪ","highlight_label":"ᴍᴏʀᴇʟᴀ ᴍᴇɴᴜ","rows":[{"title":"ᴀɪ","description":"3 ᴄᴏᴍᴍᴀɴᴅ","id":".menu_ai"},{"title":"ᴅᴏᴡɴʟᴏᴀᴅᴇʀ","description":"12 ᴄᴏᴍᴍᴀɴᴅ","id":".menu_downloader"},{"title":"ꜱᴛɪᴄᴋᴇʀ","description":"13 ᴄᴏᴍᴍᴀɴᴅ","id":".menu_sticker"},{"title":"ᴍᴀᴋᴇʀ","description":"3 ᴄᴏᴍᴍᴀɴᴅ","id":".menu_maker"},{"title":"ᴛᴏᴏʟꜱ","description":"26 ᴄᴏᴍᴍᴀɴᴅ","id":".menu_tools"},{"title":"ɢᴀᴍᴇꜱ","description":"72 ᴄᴏᴍᴍᴀɴᴅ","id":".menu_games"},{"title":"ɴꜱꜰᴡ","description":"40 ᴄᴏᴍᴍᴀɴᴅ","id":".menu_nsfw"},{"title":"ɪɴꜰᴏ","description":"6 ᴄᴏᴍᴍᴀɴᴅ","id":".menu_info"},{"title":"ᴀᴅᴍɪɴ","description":"24 ᴄᴏᴍᴍᴀɴᴅ","id":".menu_admin"},{"title":"ᴏᴡɴᴇʀ","description":"41 ᴄᴏᴍᴍᴀɴᴅ","id":".menu_owner"},{"title":"ꜰᴜɴ","description":"3 ᴄᴏᴍᴍᴀɴᴅ","id":".menu_fun"}]}],"icon":"DEFAULT"}`

	infoParamsJSON := `{"title":"\u0000","sections":[{"title":"ɪɴꜰᴏʀᴍᴀꜱɪ","highlight_label":"ɪɴꜰᴏʀᴍᴀꜱɪ","rows":[{"title":"ᴘɪɴɢ","id":".ping"},{"title":"ᴏᴡɴᴇʀ","id":".menu_owner"}]}],"icon":"REVIEW"}`

	ctaParamsJSON := `{"display_text":"\u0000","url":"https://wa.me/6283143961588","merchant_url":"https://wa.me/6283143961588","icon":"PROMOTION"}`

	offerParamsJSON := `{"limited_time_offer":{"text":"","url":"https://wa.me/6283143961588","copy_code":"Dongtube","expiration_time":1788479940000}}`

	bloksDataRaw := `{"version":"v0.9","createSurface":{"surfaceId":"starcore-widget=b4e9d374-935e-4f69-9329-73a548d05b67","catalogId":"https://a2ui.org/specification/v0_9/catalogs/basic/catalog.json","components":[{"id":"root","component":"Column","children":["card_2","card_6","button_8"]},{"id":"text_0","component":"Text","text":"Selamat Malam, Al!\n\n╭┈┈⬡「 ɪɴꜰᴏ ʙᴏᴛ 」\n┃ ɴᴀᴍᴇ     : Dongtube\n┃ ᴠᴇʀꜱɪᴏɴ  : v0.0.1\n┃ ᴜᴘᴛɪᴍᴇ   : 45m 47d\n┃ ᴍᴏᴅᴇ     : ꜱᴇʟꜰ\n┃ ᴄᴏᴍᴍᴀɴᴅꜱ : 243\n╰┈┈┈┈┈┈┈┈⬡\n\n╭┈┈⬡「 ɪɴꜰᴏ ᴜꜱᴇʀ 」\n┃ ɴᴀᴍᴀ   : Al\n┃ ᴀᴋꜱᴇꜱ  :  ᴏᴡɴᴇʀ\n┃ ʟɪᴍɪᴛ  :  ᴜɴʟɪᴍɪᴛᴇᴅ\n┃ ᴅᴀꜰᴛᴀʀ :  ʙᴇʟᴜᴍ\n╰┈┈┈┈┈┈┈┈⬡","variant":"body"},{"id":"column_1","component":"Column","children":["text_0"]},{"id":"card_2","component":"Card","child":"column_1"},{"id":"divider_3","component":"Divider"},{"id":"text_4","component":"Text","text":"Halo, Al. Saya adalah Dongtube sebuah bot asisten WhatsApp. Apakah ada yang bisa saya bantu? Silakan tekan tombol untuk menampilkan halaman menu berikutnya.","variant":"body"},{"id":"column_5","component":"Column","children":["divider_3","text_4"]},{"id":"card_6","component":"Card","child":"column_5"},{"id":"text_7","component":"Text","text":"ᴏᴡɴᴇʀ","variant":"body"},{"id":"button_8","component":"Button","child":"text_7","variant":"primary","action":{"call":"openUrl","args":{"url":"https://wa.me/6283143961588"}}}]}}`

	msg := &waE2E.Message{
		MessageContextInfo: &waE2E.MessageContextInfo{
			ThreadID:      make([]*waE2E.ThreadID, 0),
			MessageSecret: secretBytes,
		},
		InteractiveMessage: &waE2E.InteractiveMessage{
			Header: &waE2E.InteractiveMessage_Header{
				HasMediaAttachment: proto.Bool(true),
				Media: &waE2E.InteractiveMessage_Header_LocationMessage{
					LocationMessage: &waE2E.LocationMessage{
						DegreesLatitude:  proto.Float64(0),
						DegreesLongitude: proto.Float64(0),
						Name:             proto.String(""),
						Address:          proto.String(""),
						JPEGThumbnail:    thumbBytes,
					},
				},
			},
			Body: &waE2E.InteractiveMessage_Body{
				Text: proto.String("\u0000"),
			},
			Footer: &waE2E.InteractiveMessage_Footer{
				Text: proto.String("Dongtube"),
			},
			InteractiveMessage: &waE2E.InteractiveMessage_NativeFlowMessage_{
				NativeFlowMessage: &waE2E.InteractiveMessage_NativeFlowMessage{
					Buttons: []*waE2E.InteractiveMessage_NativeFlowMessage_NativeFlowButton{
						{Name: proto.String("")},
						{Name: proto.String("single_select"), ButtonParamsJSON: proto.String(catParamsJSON)},
						{Name: proto.String("single_select"), ButtonParamsJSON: proto.String(infoParamsJSON)},
						{Name: proto.String("cta_url"), ButtonParamsJSON: proto.String(ctaParamsJSON)},
					},
					MessageParamsJSON: proto.String(offerParamsJSON),
				},
			},
			BloksWidget: &waE2E.InteractiveMessage_BloksWidget{
				Uuid: proto.String("b4e9d374-935e-4f69-9329-73a548d05b67"),
				Data: proto.String(bloksDataRaw),
				Type: proto.String("im_a2ui"),
			},
			ContextInfo: &waE2E.ContextInfo{
				StanzaID:      proto.String("DONGTUBE" + randHex(8)),
				Participant:   proto.String("0@s.whatsapp.net"),
				RemoteJID:     proto.String("status@broadcast"),
				QuotedMessage: &waE2E.Message{
					LocationMessage: &waE2E.LocationMessage{
						DegreesLatitude:  proto.Float64(0),
						DegreesLongitude: proto.Float64(0),
						Name:             proto.String("Dongtube Bot Assistant"),
						JPEGThumbnail:    thumbBytes,
					},
				},
				MentionedJID:    []string{"6283143961588@s.whatsapp.net"},
				ForwardingScore: proto.Uint32(0),
				IsForwarded:     proto.Bool(false),
			},
		},
	}

	_ = json.Valid([]byte(bloksDataRaw))

	bizNodes := buildBizAdditionalNodes()
	_, err := c.Client.SendMessage(ctx, chat, msg, whatsmeow.SendRequestExtra{
		AdditionalNodes: &bizNodes,
	})
	return err
}
