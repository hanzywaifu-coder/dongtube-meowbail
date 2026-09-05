package meowbail

import (
	"context"
	"encoding/json"
	"fmt"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// SendExactDongtubeMenu mengirim menu interaktif tombol persis referensi 7e716445aee7.bin:
// Memasukkan teks menu ke Body agar 100% terbaca di semua WhatsApp (Android, iOS, Web, Desktop),
// plus 3 tombol native flow (single_select kategori, single_select info, cta_url owner)
// dan locationMessage header dengan jpegThumbnail
func (c *Client) SendExactDongtubeMenu(ctx context.Context, chat types.JID, thumbData []byte, botName, ownerPhone string, uptimeStr string, cmdCount int, userName string, greeting string) error {
	if userName == "" {
		userName = "Al"
	}
	if botName == "" {
		botName = "Dongtube"
	}
	if ownerPhone == "" {
		ownerPhone = "6283143961588"
	}
	if greeting == "" {
		greeting = "Selamat Malam"
	}

	// 1. Baris Kategori Dropdown 1
	categoryRows := []map[string]interface{}{
		{"title": "ᴀɪ", "description": "3 ᴄᴏᴍᴍᴀɴᴅ", "id": ".menu_ai"},
		{"title": "ᴅᴏᴡɴʟᴏᴀᴅᴇʀ", "description": "12 ᴄᴏᴍᴍᴀɴᴅ", "id": ".menu_downloader"},
		{"title": "ꜱᴛɪᴄᴋᴇʀ", "description": "13 ᴄᴏᴍᴍᴀɴᴅ", "id": ".menu_sticker"},
		{"title": "ᴍᴀᴋᴇʀ", "description": "3 ᴄᴏᴍᴍᴀɴᴅ", "id": ".menu_maker"},
		{"title": "ᴛᴏᴏʟꜱ", "description": "26 ᴄᴏᴍᴍᴀɴᴅ", "id": ".menu_tools"},
		{"title": "ɢᴀᴍᴇꜱ", "description": "72 ᴄᴏᴍᴍᴀɴᴅ", "id": ".menu_games"},
		{"title": "ɴꜱꜰᴡ", "description": "40 ᴄᴏᴍᴍᴀɴᴅ", "id": ".menu_nsfw"},
		{"title": "ɪɴꜰᴏ", "description": "6 ᴄᴏᴍᴍᴀɴᴅ", "id": ".menu_info"},
		{"title": "ᴀᴅᴍɪɴ", "description": "24 ᴄᴏᴍᴍᴀɴᴅ", "id": ".menu_admin"},
		{"title": "ᴏᴡɴᴇʀ", "description": "41 ᴄᴏᴍᴍᴀɴᴅ", "id": ".menu_owner"},
		{"title": "ꜰᴜɴ", "description": "3 ᴄᴏᴍᴍᴀɴᴅ", "id": ".menu_fun"},
	}

	catParamsJSON, _ := json.Marshal(map[string]interface{}{
		"title": "ᴋᴀᴛᴇɢᴏʀɪ",
		"sections": []map[string]interface{}{
			{
				"title":           "ᴋᴀᴛᴇɢᴏʀɪ",
				"highlight_label": "ᴍᴏʀᴇʟᴀ ᴍᴇɴᴜ",
				"rows":            categoryRows,
			},
		},
		"icon": "DEFAULT",
	})

	// 2. Baris Informasi Dropdown 2
	infoParamsJSON, _ := json.Marshal(map[string]interface{}{
		"title": "ɪɴꜰᴏʀᴍᴀꜱɪ",
		"sections": []map[string]interface{}{
			{
				"title":           "ɪɴꜰᴏʀᴍᴀꜱɪ",
				"highlight_label": "ɪɴꜰᴏʀᴍᴀꜱɪ",
				"rows": []map[string]interface{}{
					{"title": "ᴘɪɴɢ", "description": "Cek kecepatan bot", "id": ".ping"},
					{"title": "ᴏᴡɴᴇʀ", "description": "Menu khusus owner", "id": ".menu_owner"},
				},
			},
		},
		"icon": "REVIEW",
	})

	// 3. Tombol CTA URL
	ctaParamsJSON, _ := json.Marshal(map[string]interface{}{
		"display_text": "ᴏᴡɴᴇʀ",
		"url":          fmt.Sprintf("https://wa.me/%s", ownerPhone),
		"merchant_url": fmt.Sprintf("https://wa.me/%s", ownerPhone),
		"icon":         "PROMOTION",
	})

	offerParamsJSON, _ := json.Marshal(map[string]interface{}{
		"limited_time_offer": map[string]interface{}{
			"text":            "",
			"url":             fmt.Sprintf("https://wa.me/%s", ownerPhone),
			"copy_code":       botName,
			"expiration_time": 1788479940000,
		},
	})

	nativeButtons := []*waE2E.InteractiveMessage_NativeFlowMessage_NativeFlowButton{
		{Name: proto.String("single_select"), ButtonParamsJSON: proto.String(string(catParamsJSON))},
		{Name: proto.String("single_select"), ButtonParamsJSON: proto.String(string(infoParamsJSON))},
		{Name: proto.String("cta_url"), ButtonParamsJSON: proto.String(string(ctaParamsJSON))},
	}

	// Teks Menu Tampil di Layar Chat
	displayText := fmt.Sprintf("%s, %s!\n\n"+
		"╭┈┈⬡「 ɪɴꜰᴏ ʙᴏᴛ 」\n"+
		"┃ ɴᴀᴍᴇ     : %s\n"+
		"┃ ᴠᴇʀꜱɪᴏɴ  : v0.0.1\n"+
		"┃ ᴜᴘᴛɪᴍᴇ   : %s\n"+
		"┃ ᴍᴏᴅᴇ     : ꜱᴇʟꜰ\n"+
		"┃ ᴄᴏᴍᴍᴀɴᴅꜱ : %d\n"+
		"╰┈┈┈┈┈┈┈┈⬡\n\n"+
		"╭┈┈⬡「 ɪɴꜰᴏ ᴜꜱᴇʀ 」\n"+
		"┃ ɴᴀᴍᴀ   : %s\n"+
		"┃ ᴀᴋꜱᴇꜱ  :  ᴏᴡɴᴇʀ\n"+
		"┃ ʟɪᴍɪᴛ  :  ᴜɴʟɪᴍɪᴛᴇᴅ\n"+
		"┃ ᴅᴀꜰᴛᴀʀ :  ʙᴇʟᴜᴍ\n"+
		"╰┈┈┈┈┈┈┈┈⬡\n\n"+
		"Halo, %s. Saya adalah %s sebuah bot asisten WhatsApp. Apakah ada yang bisa saya bantu? Silakan tekan tombol untuk menampilkan halaman menu berikutnya.",
		greeting, userName, botName, uptimeStr, cmdCount, userName, userName, botName)

	// Header LocationMessage
	header := &waE2E.InteractiveMessage_Header{
		HasMediaAttachment: proto.Bool(len(thumbData) > 0),
		Media: &waE2E.InteractiveMessage_Header_LocationMessage{
			LocationMessage: &waE2E.LocationMessage{
				DegreesLatitude:  proto.Float64(0),
				DegreesLongitude: proto.Float64(0),
				Name:             proto.String(""),
				Address:          proto.String(""),
				JPEGThumbnail:    thumbData,
			},
		},
	}

	interactiveMsg := &waE2E.InteractiveMessage{
		Header: header,
		Body: &waE2E.InteractiveMessage_Body{
			Text: proto.String(displayText),
		},
		Footer: &waE2E.InteractiveMessage_Footer{
			Text: proto.String(botName),
		},
		InteractiveMessage: &waE2E.InteractiveMessage_NativeFlowMessage_{
			NativeFlowMessage: &waE2E.InteractiveMessage_NativeFlowMessage{
				Buttons:           nativeButtons,
				MessageParamsJSON: proto.String(string(offerParamsJSON)),
			},
		},
		ContextInfo: &waE2E.ContextInfo{
			MentionedJID:    []string{fmt.Sprintf("%s@s.whatsapp.net", ownerPhone)},
			ForwardingScore: proto.Uint32(0),
			IsForwarded:     proto.Bool(false),
		},
	}

	msg := &waE2E.Message{
		ViewOnceMessage: &waE2E.FutureProofMessage{
			Message: &waE2E.Message{
				InteractiveMessage: interactiveMsg,
			},
		},
	}

	bizNodes := buildBizAdditionalNodes()
	_, err := c.Client.SendMessage(ctx, chat, msg, whatsmeow.SendRequestExtra{
		AdditionalNodes: &bizNodes,
	})
	return err
}
