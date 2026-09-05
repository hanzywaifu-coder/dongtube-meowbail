package meowbail

import (
	"context"
	"encoding/json"
	"fmt"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waCommon"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

// SendText sends a plain text message
func (c *Client) SendText(ctx context.Context, chat types.JID, text string, opts ...*MessageOptions) error {
	msg := &waE2E.Message{
		Conversation: &text,
	}

	if len(opts) > 0 && opts[0] != nil {
		msg = applyOptions(msg, opts[0])
	}

	_, err := c.Client.SendMessage(ctx, chat, msg)
	return err
}

// SendTextWithNewsletter sends text with newsletter context
func (c *Client) SendTextWithNewsletter(ctx context.Context, chat types.JID, text string) error {
	msg := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(text),
			ContextInfo: buildNewsletterContext(c.config),
		},
	}

	_, err := c.Client.SendMessage(ctx, chat, msg)
	return err
}

// SendButtons sends a message with quick_reply buttons
func (c *Client) SendButtons(ctx context.Context, chat types.JID, text string, buttons []Button, opts ...*MessageOptions) error {
	if len(buttons) == 0 {
		return c.SendText(ctx, chat, text, opts...)
	}

	// Build ButtonsMessage
	var waButtons []*waE2E.ButtonsMessage_Button
	for i, btn := range buttons {
		switch btn.Type {
		case ButtonQuickReply:
			waButtons = append(waButtons, &waE2E.ButtonsMessage_Button{
				ButtonID: proto.String(btn.ID),
				ButtonText: &waE2E.ButtonsMessage_Button_ButtonText{
					DisplayText: proto.String(btn.Text),
				},
				Type: waE2E.ButtonsMessage_Button_RESPONSE.Enum(),
			})
		case ButtonCTAURL:
			params, _ := json.Marshal(map[string]interface{}{
				"display_text": btn.DisplayText,
				"url":          btn.URL,
				"merchant_url": btn.URL,
			})
			waButtons = append(waButtons, &waE2E.ButtonsMessage_Button{
				ButtonID: proto.String(fmt.Sprintf("cta_%d", i)),
				ButtonText: &waE2E.ButtonsMessage_Button_ButtonText{
					DisplayText: proto.String(btn.Text),
				},
				Type: waE2E.ButtonsMessage_Button_NATIVE_FLOW.Enum(),
				NativeFlowInfo: &waE2E.ButtonsMessage_Button_NativeFlowInfo{
					Name:       proto.String("cta_url"),
					ParamsJSON: proto.String(string(params)),
				},
			})
		case ButtonPhoneNumber:
			params, _ := json.Marshal(map[string]interface{}{
				"display_text": btn.DisplayText,
				"phone_number": btn.Phone,
			})
			waButtons = append(waButtons, &waE2E.ButtonsMessage_Button{
				ButtonID: proto.String(fmt.Sprintf("phone_%d", i)),
				ButtonText: &waE2E.ButtonsMessage_Button_ButtonText{
					DisplayText: proto.String(btn.Text),
				},
				Type: waE2E.ButtonsMessage_Button_NATIVE_FLOW.Enum(),
				NativeFlowInfo: &waE2E.ButtonsMessage_Button_NativeFlowInfo{
					Name:       proto.String("quick_reply"),
					ParamsJSON: proto.String(string(params)),
				},
			})
		case ButtonCopyText:
			params, _ := json.Marshal(map[string]interface{}{
				"display_text": btn.DisplayText,
				"copy_code":    btn.ID,
			})
			waButtons = append(waButtons, &waE2E.ButtonsMessage_Button{
				ButtonID: proto.String(fmt.Sprintf("copy_%d", i)),
				ButtonText: &waE2E.ButtonsMessage_Button_ButtonText{
					DisplayText: proto.String(btn.Text),
				},
				Type: waE2E.ButtonsMessage_Button_NATIVE_FLOW.Enum(),
				NativeFlowInfo: &waE2E.ButtonsMessage_Button_NativeFlowInfo{
					Name:       proto.String("quick_reply"),
					ParamsJSON: proto.String(string(params)),
				},
			})
		}
	}

	msg := &waE2E.Message{
		ButtonsMessage: &waE2E.ButtonsMessage{
			HeaderType: waE2E.ButtonsMessage_TEXT.Enum(),
			Header: &waE2E.ButtonsMessage_Text{
				Text: "",
			},
			ContentText: proto.String(text),
			FooterText:  proto.String(""),
			Buttons:     waButtons,
		},
	}

	_, err := c.Client.SendMessage(ctx, chat, msg)
	return err
}

// SendList sends a list/dropdown menu message
func (c *Client) SendList(ctx context.Context, chat types.JID, title, description, buttonText string, sections []Section, opts ...*MessageOptions) error {
	var waSections []*waE2E.ListMessage_Section
	for _, sec := range sections {
		var rows []*waE2E.ListMessage_Row
		for _, row := range sec.Rows {
			rows = append(rows, &waE2E.ListMessage_Row{
				Title:       proto.String(row.Title),
				Description: proto.String(row.Description),
				RowID:       proto.String(row.ID),
			})
		}
		waSections = append(waSections, &waE2E.ListMessage_Section{
			Title: proto.String(sec.Title),
			Rows:  rows,
		})
	}

	msg := &waE2E.Message{
		ListMessage: &waE2E.ListMessage{
			Title:      proto.String(title),
			Description: proto.String(description),
			ButtonText: proto.String(buttonText),
			ListType:   waE2E.ListMessage_SINGLE_SELECT.Enum(),
			Sections:   waSections,
			FooterText: proto.String(""),
			ContextInfo: buildNewsletterContext(c.config),
		},
	}

	_, err := c.Client.SendMessage(ctx, chat, msg)
	return err
}

// SendImage sends an image message
func (c *Client) SendImage(ctx context.Context, chat types.JID, data []byte, caption string, opts ...*MessageOptions) error {
	resp, err := c.Client.Upload(ctx, data, whatsmeow.MediaImage)
	if err != nil {
		return err
	}

	msg := &waE2E.Message{
		ImageMessage: &waE2E.ImageMessage{
			URL:           &resp.URL,
			Mimetype:      proto.String("image/jpeg"),
			Caption:       proto.String(caption),
			FileEncSHA256: resp.FileEncSHA256,
			FileSHA256:    resp.FileSHA256,
			FileLength:    &resp.FileLength,
			DirectPath:    &resp.DirectPath,
			MediaKey:      resp.MediaKey,
			ContextInfo:   buildNewsletterContext(c.config),
		},
	}

	_, err = c.Client.SendMessage(ctx, chat, msg)
	return err
}

// SendDocument sends a document message
func (c *Client) SendDocument(ctx context.Context, chat types.JID, data []byte, filename, mimetype, caption string, opts ...*MessageOptions) error {
	resp, err := c.Client.Upload(ctx, data, whatsmeow.MediaDocument)
	if err != nil {
		return err
	}

	msg := &waE2E.Message{
		DocumentMessage: &waE2E.DocumentMessage{
			URL:           &resp.URL,
			Mimetype:      proto.String(mimetype),
			FileName:      proto.String(filename),
			Caption:       proto.String(caption),
			FileEncSHA256: resp.FileEncSHA256,
			FileSHA256:    resp.FileSHA256,
			FileLength:    &resp.FileLength,
			DirectPath:    &resp.DirectPath,
			MediaKey:      resp.MediaKey,
			ContextInfo:   buildNewsletterContext(c.config),
		},
	}

	_, err = c.Client.SendMessage(ctx, chat, msg)
	return err
}

// SendVideo sends a video message
func (c *Client) SendVideo(ctx context.Context, chat types.JID, data []byte, caption string, opts ...*MessageOptions) error {
	resp, err := c.Client.Upload(ctx, data, whatsmeow.MediaVideo)
	if err != nil {
		return err
	}

	msg := &waE2E.Message{
		VideoMessage: &waE2E.VideoMessage{
			URL:           &resp.URL,
			Mimetype:      proto.String("video/mp4"),
			Caption:       proto.String(caption),
			FileEncSHA256: resp.FileEncSHA256,
			FileSHA256:    resp.FileSHA256,
			FileLength:    &resp.FileLength,
			DirectPath:    &resp.DirectPath,
			MediaKey:      resp.MediaKey,
			ContextInfo:   buildNewsletterContext(c.config),
		},
	}

	_, err = c.Client.SendMessage(ctx, chat, msg)
	return err
}

// SendAudio sends an audio message
func (c *Client) SendAudio(ctx context.Context, chat types.JID, data []byte, opts ...*MessageOptions) error {
	resp, err := c.Client.Upload(ctx, data, whatsmeow.MediaAudio)
	if err != nil {
		return err
	}

	msg := &waE2E.Message{
		AudioMessage: &waE2E.AudioMessage{
			URL:           &resp.URL,
			Mimetype:      proto.String("audio/mp4"),
			FileEncSHA256: resp.FileEncSHA256,
			FileSHA256:    resp.FileSHA256,
			FileLength:    &resp.FileLength,
			DirectPath:    &resp.DirectPath,
			MediaKey:      resp.MediaKey,
		},
	}

	_, err = c.Client.SendMessage(ctx, chat, msg)
	return err
}

// SendSticker sends a sticker message
func (c *Client) SendSticker(ctx context.Context, chat types.JID, data []byte, opts ...*MessageOptions) error {
	resp, err := c.Client.Upload(ctx, data, whatsmeow.MediaImage)
	if err != nil {
		return err
	}

	msg := &waE2E.Message{
		StickerMessage: &waE2E.StickerMessage{
			URL:           &resp.URL,
			Mimetype:      proto.String("image/webp"),
			FileEncSHA256: resp.FileEncSHA256,
			FileSHA256:    resp.FileSHA256,
			FileLength:    &resp.FileLength,
			DirectPath:    &resp.DirectPath,
			MediaKey:      resp.MediaKey,
		},
	}

	_, err = c.Client.SendMessage(ctx, chat, msg)
	return err
}

// SendLocation sends a location message
func (c *Client) SendLocation(ctx context.Context, chat types.JID, lat, lng float64, name, address string) error {
	msg := &waE2E.Message{
		LocationMessage: &waE2E.LocationMessage{
			DegreesLatitude:  proto.Float64(lat),
			DegreesLongitude: proto.Float64(lng),
			Name:             proto.String(name),
			Address:          proto.String(address),
		},
	}

	_, err := c.Client.SendMessage(ctx, chat, msg)
	return err
}

// SendContact sends a contact/vcard message
func (c *Client) SendContact(ctx context.Context, chat types.JID, name, phone string) error {
	vcard := fmt.Sprintf("BEGIN:VCARD\nVERSION:3.0\nFN:%s\nTEL;type=CELL;type=VOICE;waid=%s:+%s\nEND:VCARD", name, phone, phone)

	msg := &waE2E.Message{
		ContactMessage: &waE2E.ContactMessage{
			DisplayName: proto.String(name),
			Vcard:       proto.String(vcard),
		},
	}

	_, err := c.Client.SendMessage(ctx, chat, msg)
	return err
}

// SendReaction sends a reaction to a message
func (c *Client) SendReaction(ctx context.Context, chat types.JID, msgID types.MessageID, emoji string) error {
	reaction := &waE2E.ReactionMessage{
		Key: &waCommon.MessageKey{
			RemoteJID: proto.String(chat.String()),
			ID:        proto.String(string(msgID)),
		},
		Text: proto.String(emoji),
	}

	_, err := c.Client.SendMessage(ctx, chat, &waE2E.Message{
		ReactionMessage: reaction,
	})
	return err
}

// SendPoll sends a poll message
func (c *Client) SendPoll(ctx context.Context, chat types.JID, name string, options []string, selectableCount int) error {
	var pollOptions []*waE2E.PollCreationMessage_Option
	for _, opt := range options {
		pollOptions = append(pollOptions, &waE2E.PollCreationMessage_Option{
			OptionName: proto.String(opt),
		})
	}

	msg := &waE2E.Message{
		PollCreationMessage: &waE2E.PollCreationMessage{
			Name:                   proto.String(name),
			Options:                pollOptions,
			SelectableOptionsCount: proto.Uint32(uint32(selectableCount)),
		},
	}

	_, err := c.Client.SendMessage(ctx, chat, msg)
	return err
}

// HandleButtonResponse checks if a message event is a button response
func HandleButtonResponse(evt interface{}) (string, bool) {
	e, ok := evt.(*events.Message)
	if !ok {
		return "", false
	}

	if e.Message.ButtonsResponseMessage != nil {
		selectedID := e.Message.ButtonsResponseMessage.GetSelectedButtonID()
		if selectedID != "" {
			return selectedID, true
		}
	}

	if e.Message.ListResponseMessage != nil {
		singleSelect := e.Message.ListResponseMessage.GetSingleSelectReply()
		if singleSelect != nil {
			return singleSelect.GetSelectedRowID(), true
		}
	}

	return "", false
}

// Helper functions

func buildNewsletterContext(cfg *Config) *waE2E.ContextInfo {
	if cfg == nil || cfg.NewsletterJID == "" {
		return &waE2E.ContextInfo{}
	}

	return &waE2E.ContextInfo{
		IsForwarded:     proto.Bool(true),
		ForwardingScore: proto.Uint32(9999),
		BusinessMessageForwardInfo: &waE2E.ContextInfo_BusinessMessageForwardInfo{
			BusinessOwnerJID: proto.String(cfg.BusinessOwnerJID),
		},
		ForwardedNewsletterMessageInfo: &waE2E.ContextInfo_ForwardedNewsletterMessageInfo{
			NewsletterJID:  proto.String(cfg.NewsletterJID),
			NewsletterName: proto.String(cfg.NewsletterName),
		},
	}
}

func applyOptions(msg *waE2E.Message, opts *MessageOptions) *waE2E.Message {
	if opts == nil {
		return msg
	}

	// Apply mentions
	if len(opts.Mentions) > 0 {
		setMentions(msg, opts.Mentions)
	}

	return msg
}

func setMentions(msg *waE2E.Message, mentions []string) {
	if msg.ExtendedTextMessage != nil {
		msg.ExtendedTextMessage.ContextInfo = &waE2E.ContextInfo{
			MentionedJID: mentions,
		}
	} else if msg.Conversation != nil {
		text := *msg.Conversation
		msg.ExtendedTextMessage = &waE2E.ExtendedTextMessage{
			Text: proto.String(text),
			ContextInfo: &waE2E.ContextInfo{
				MentionedJID: mentions,
			},
		}
		msg.Conversation = nil
	}
}
