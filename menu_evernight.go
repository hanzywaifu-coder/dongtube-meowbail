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

// SendEvernightMenu mengirim menu persis Baileys Evernight AI:
// - Header: DocumentMessage (package.json dummy) + jpegThumbnail + pageCount 100
// - Footer: menuBody teks lengkap info bot & user
// - NativeFlow buttons: single_select dropdown & cta_url info saluran
// - ViewOnceMessage wrapper
// - AdditionalNodes: <biz> stanza
func (c *Client) SendEvernightMenu(ctx context.Context, chat types.JID, docData []byte, thumbData []byte, botName, menuBodyText string, sections []Section, ctaText, ctaURL string) error {
	// 1. Upload media document jika belum
	var docMsg *waE2E.DocumentMessage
	if len(docData) > 0 {
		resp, err := c.Client.Upload(ctx, docData, whatsmeow.MediaDocument)
		if err == nil {
			docLen := uint64(len(docData))
			docMsg = &waE2E.DocumentMessage{
				URL:           &resp.URL,
				Mimetype:      proto.String("image/png"),
				FileName:      proto.String(botName),
				FileLength:    &docLen,
				PageCount:     proto.Uint32(100),
				FileEncSHA256: resp.FileEncSHA256,
				FileSHA256:    resp.FileSHA256,
				DirectPath:    &resp.DirectPath,
				MediaKey:      resp.MediaKey,
				JPEGThumbnail: thumbData,
			}
		}
	}

	// 2. Format Sections Dropdown
	var waSections []map[string]interface{}
	for _, sec := range sections {
		var rows []map[string]interface{}
		for _, r := range sec.Rows {
			rows = append(rows, map[string]interface{}{
				"title":       r.Title,
				"description": r.Description,
				"id":          r.ID,
			})
		}
		waSections = append(waSections, map[string]interface{}{
			"title": sec.Title,
			"rows":  rows,
		})
	}

	catParamsJSON, _ := json.Marshal(map[string]interface{}{
		"title":    "Selection",
		"sections": waSections,
	})

	nativeButtons := []*waE2E.InteractiveMessage_NativeFlowMessage_NativeFlowButton{
		{
			Name:             proto.String("single_select"),
			ButtonParamsJSON: proto.String(string(catParamsJSON)),
		},
	}

	if ctaURL != "" {
		if ctaText == "" {
			ctaText = fmt.Sprintf("Info %s", botName)
		}
		ctaParamsJSON, _ := json.Marshal(map[string]interface{}{
			"display_text": ctaText,
			"url":          ctaURL,
			"merchant_url": ctaURL,
		})
		nativeButtons = append(nativeButtons, &waE2E.InteractiveMessage_NativeFlowMessage_NativeFlowButton{
			Name:             proto.String("cta_url"),
			ButtonParamsJSON: proto.String(string(ctaParamsJSON)),
		})
	}

	// 3. Header Document Message
	header := &waE2E.InteractiveMessage_Header{
		Title:              proto.String(""),
		HasMediaAttachment: proto.Bool(docMsg != nil),
	}
	if docMsg != nil {
		header.Media = &waE2E.InteractiveMessage_Header_DocumentMessage{
			DocumentMessage: docMsg,
		}
	}

	interactiveMsg := &waE2E.InteractiveMessage{
		Header: header,
		Body: &waE2E.InteractiveMessage_Body{
			Text: proto.String(""),
		},
		Footer: &waE2E.InteractiveMessage_Footer{
			Text: proto.String(menuBodyText),
		},
		InteractiveMessage: &waE2E.InteractiveMessage_NativeFlowMessage_{
			NativeFlowMessage: &waE2E.InteractiveMessage_NativeFlowMessage{
				Buttons: nativeButtons,
			},
		},
		ContextInfo: buildNewsletterContext(c.config),
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
