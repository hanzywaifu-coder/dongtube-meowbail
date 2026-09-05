package meowbail

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"go.mau.fi/whatsmeow"
	waBinary "go.mau.fi/whatsmeow/binary"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// buildBizAdditionalNodes membangun XML stanza node <biz> yang diwajibkan oleh protokol WhatsApp
// untuk merender tombol interaktif (native_flow, single_select, cta_url, dll)
func buildBizAdditionalNodes() []waBinary.Node {
	ts := strconv.FormatInt(time.Now().Unix()-77980457, 10)
	return []waBinary.Node{
		{
			Tag: "biz",
			Attrs: waBinary.Attrs{
				"actual_actors":   "2",
				"host_storage":    "2",
				"privacy_mode_ts": ts,
			},
			Content: []waBinary.Node{
				{
					Tag: "engagement",
					Attrs: waBinary.Attrs{
						"customer_service_state": "open",
						"conversation_state":     "open",
					},
				},
				{
					Tag: "interactive",
					Attrs: waBinary.Attrs{
						"type": "native_flow",
						"v":    "1",
					},
					Content: []waBinary.Node{
						{
							Tag: "native_flow",
							Attrs: waBinary.Attrs{
								"name": "mixed",
								"v":    "9",
							},
						},
					},
				},
			},
		},
	}
}

// SendExactA2UIMenu mengirim format menu exact A2UI BloksWidget + NativeFlow Buttons
// lengkap dengan stanza XML biz node injection agar tombol benar-benar muncul di WhatsApp!
func (c *Client) SendExactA2UIMenu(ctx context.Context, chat types.JID, thumbData []byte, botName string, ownerPhone string, sections []Section, uptimeStr string, cmdCount int, userName string) error {
	if userName == "" {
		userName = "User"
	}
	if botName == "" {
		botName = "Dongtube"
	}
	if ownerPhone == "" {
		ownerPhone = "6283143961588"
	}

	// 1. Format sections kategori
	var catSections []map[string]interface{}
	for _, sec := range sections {
		var rows []map[string]interface{}
		for _, r := range sec.Rows {
			rows = append(rows, map[string]interface{}{
				"title":       r.Title,
				"description": r.Description,
				"id":          r.ID,
			})
		}
		catSections = append(catSections, map[string]interface{}{
			"title":           sec.Title,
			"highlight_label": "MORELA MENU",
			"rows":            rows,
		})
	}

	catParamsJSON, _ := json.Marshal(map[string]interface{}{
		"title":    "\x00",
		"sections": catSections,
		"icon":     "DEFAULT",
	})

	infoParamsJSON, _ := json.Marshal(map[string]interface{}{
		"title": "\x00",
		"sections": []map[string]interface{}{
			{
				"title":           "ɪɴꜰᴏʀᴍᴀꜱɪ",
				"highlight_label": "ɪɴꜰᴏʀᴍᴀꜱɪ",
				"rows": []map[string]interface{}{
					{"title": "ᴘɪɴɢ", "id": ".ping"},
					{"title": "ᴏᴡɴᴇʀ", "id": ".owner"},
				},
			},
		},
		"icon": "REVIEW",
	})

	ctaParamsJSON, _ := json.Marshal(map[string]interface{}{
		"display_text": "\x00",
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

	// 2. Data BloksWidget UI
	bloksText0 := fmt.Sprintf("Selamat Datang, %s!\n\n"+
		"╭┈┈⬡「 ɪɴꜰᴏ ʙᴏᴛ 」\n"+
		"┃ ɴᴀᴍᴇ     : %s\n"+
		"┃ ᴠᴇʀꜱɪᴏɴ  : v1.0.0\n"+
		"┃ ᴜᴘᴛɪᴍᴇ   : %s\n"+
		"┃ ᴍᴏᴅᴇ     : ꜱᴇʟꜰ\n"+
		"┃ ᴄᴏᴍᴍᴀɴᴅꜱ : %d\n"+
		"╰┈┈┈┈┈┈┈┈⬡\n\n"+
		"╭┈┈⬡「 ɪɴꜰᴏ ᴜꜱᴇʀ 」\n"+
		"┃ ɴᴀᴍᴀ   : %s\n"+
		"┃ ᴀᴋꜱᴇꜱ  : ᴏᴡɴᴇʀ\n"+
		"┃ ʟɪᴍɪᴛ  : ᴜɴʟɪᴍɪᴛᴇᴅ\n"+
		"┃ ᴅᴀꜰᴛᴀʀ : ꜱᴜᴅᴀʜ\n"+
		"╰┈┈┈┈┈┈┈┈⬡", userName, botName, uptimeStr, cmdCount, userName)

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

	// 3. Header dengan LocationMessage & JPEGThumbnail persis referensi
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
			Text: proto.String("\x00"),
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

	// KUNCI UTAMA: Injeksi AdditionalNodes <biz> ke stanza XML pengiriman pesan
	bizNodes := buildBizAdditionalNodes()
	_, err := c.Client.SendMessage(ctx, chat, msg, whatsmeow.SendRequestExtra{
		AdditionalNodes: &bizNodes,
	})
	return err
}
