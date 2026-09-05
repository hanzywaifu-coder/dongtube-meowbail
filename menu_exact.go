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

// SendExactDongtubeMenu mengirim menu interaktif persis 100% spesifikasi referensi:
// - Header: LocationMessage (lat:0, long:0, jpegThumbnail)
// - Body: "\u0000"
// - Footer: "Dongtube"
// - NativeFlow buttons: 4 buttons (empty, single_select kategori, single_select info, cta_url)
// - limited_time_offer metadata
// - bloksWidget type "im_a2ui" dengan Card, Divider, Button openUrl
// - Injeksi <biz> XML node
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

	// 1. Buttons List persis referensi
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
		"title": "\u0000",
		"sections": []map[string]interface{}{
			{
				"title":           "ᴋᴀᴛᴇɢᴏʀɪ",
				"highlight_label": "ᴍᴏʀᴇʟᴀ ᴍᴇɴᴜ",
				"rows":            categoryRows,
			},
		},
		"icon": "DEFAULT",
	})

	infoParamsJSON, _ := json.Marshal(map[string]interface{}{
		"title": "\u0000",
		"sections": []map[string]interface{}{
			{
				"title":           "ɪɴꜰᴏʀᴍᴀꜱɪ",
				"highlight_label": "ɪɴꜰᴏʀᴍᴀꜱɪ",
				"rows": []map[string]interface{}{
					{"title": "ᴘɪɴɢ", "id": ".ping"},
					{"title": "ᴏᴡɴᴇʀ", "id": ".menu_owner"},
				},
			},
		},
		"icon": "REVIEW",
	})

	ctaParamsJSON, _ := json.Marshal(map[string]interface{}{
		"display_text": "\u0000",
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
		{Name: proto.String("")},
		{Name: proto.String("single_select"), ButtonParamsJSON: proto.String(string(catParamsJSON))},
		{Name: proto.String("single_select"), ButtonParamsJSON: proto.String(string(infoParamsJSON))},
		{Name: proto.String("cta_url"), ButtonParamsJSON: proto.String(string(ctaParamsJSON))},
	}

	// 2. BloksWidget A2UI JSON Data (Persis 100% referensi)
	bloksText0 := fmt.Sprintf("%s, %s!\n\n"+
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
		"╰┈┈┈┈┈┈┈┈⬡", greeting, userName, botName, uptimeStr, cmdCount, userName)

	bloksText4 := fmt.Sprintf("Halo, %s. Saya adalah %s sebuah bot asisten WhatsApp. Apakah ada yang bisa saya bantu? Silakan tekan tombol untuk menampilkan halaman menu berikutnya.", userName, botName)

	bloksDataMap := map[string]interface{}{
		"version": "v0.9",
		"createSurface": map[string]interface{}{
			"surfaceId": "starcore-widget=b4e9d374-935e-4f69-9329-73a548d05b67",
			"catalogId": "https://a2ui.org/specification/v0_9/catalogs/basic/catalog.json",
			"components": []map[string]interface{}{
				{"id": "root", "component": "Column", "children": []string{"card_2", "card_6", "button_8"}},
				{"id": "text_0", "component": "Text", "text": bloksText0, "variant": "body"},
				{"id": "column_1", "component": "Column", "children": []string{"text_0"}},
				{"id": "card_2", "component": "Card", "child": "column_1"},
				{"id": "divider_3", "component": "Divider"},
				{"id": "text_4", "component": "Text", "text": bloksText4, "variant": "body"},
				{"id": "column_5", "component": "Column", "children": []string{"divider_3", "text_4"}},
				{"id": "card_6", "component": "Card", "child": "column_5"},
				{"id": "text_7", "component": "Text", "text": "ᴏᴡɴᴇʀ", "variant": "body"},
				{
					"id": "button_8", "component": "Button", "child": "text_7", "variant": "primary",
					"action": map[string]interface{}{
						"call": "openUrl",
						"args": map[string]interface{}{"url": fmt.Sprintf("https://wa.me/%s", ownerPhone)},
					},
				},
			},
		},
	}
	bloksJSONBytes, _ := json.Marshal(bloksDataMap)

	// 3. Header LocationMessage
	header := &waE2E.InteractiveMessage_Header{
		HasMediaAttachment: proto.Bool(true),
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
			Text: proto.String("\u0000"),
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
		BloksWidget: &waE2E.InteractiveMessage_BloksWidget{
			Uuid: proto.String("b4e9d374-935e-4f69-9329-73a548d05b67"),
			Data: proto.String(string(bloksJSONBytes)),
			Type: proto.String("im_a2ui"),
		},
		ContextInfo: &waE2E.ContextInfo{
			MentionedJID:    []string{fmt.Sprintf("%s@s.whatsapp.net", ownerPhone)},
			ForwardingScore: proto.Uint32(0),
			IsForwarded:     proto.Bool(false),
		},
	}

	msg := &waE2E.Message{
		MessageContextInfo: &waE2E.MessageContextInfo{
			MessageSecret: []byte("IXM8c2x6FioXJZLibSUFkhheds8R4KQtoWqKWvhwIkY="),
		},
		InteractiveMessage: interactiveMsg,
	}

	bizNodes := buildBizAdditionalNodes()
	_, err := c.Client.SendMessage(ctx, chat, msg, whatsmeow.SendRequestExtra{
		AdditionalNodes: &bizNodes,
	})
	return err
}
