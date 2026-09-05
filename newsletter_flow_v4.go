package meowbail

import (
	"context"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// NewsletterReactionMode mode reaksi emoji di saluran (ALL, BASIC, NONE)
type NewsletterReactionMode string

const (
	NewsletterReactionAll   NewsletterReactionMode = "ALL"
	NewsletterReactionBasic NewsletterReactionMode = "BASIC"
	NewsletterReactionBlock NewsletterReactionMode = "BLOCK"
	NewsletterReactionNone  NewsletterReactionMode = "NONE"
)

// SendNewsletterPoll membuat polling di saluran/newsletter WhatsApp
func (c *Client) SendNewsletterPoll(ctx context.Context, channelJID types.JID, name string, options []string, selectableCount int) error {
	if selectableCount <= 0 {
		selectableCount = 1
	}

	optList := make([]*waE2E.PollCreationMessage_Option, len(options))
	for i, opt := range options {
		optList[i] = &waE2E.PollCreationMessage_Option{
			OptionName: proto.String(opt),
		}
	}

	msg := &waE2E.Message{
		PollCreationMessageV3: &waE2E.PollCreationMessage{
			Name:                   proto.String(name),
			Options:                optList,
			SelectableOptionsCount: proto.Uint32(uint32(selectableCount)),
		},
	}

	_, err := c.Client.SendMessage(ctx, channelJID, msg)
	return err
}

// SendNewsletterVideo mengirim video ke saluran/newsletter
func (c *Client) SendNewsletterVideo(ctx context.Context, channelJID types.JID, videoBytes []byte, caption string) error {
	return c.SendVideo(ctx, channelJID, videoBytes, caption)
}

// SendNewsletterAudio mengirim rekaman audio/suara ke saluran
func (c *Client) SendNewsletterAudio(ctx context.Context, channelJID types.JID, audioBytes []byte) error {
	return c.SendAudio(ctx, channelJID, audioBytes)
}

// SendNewsletterDocument mengirim file dokumen ke saluran
func (c *Client) SendNewsletterDocument(ctx context.Context, channelJID types.JID, docBytes []byte, filename, mimetype, caption string) error {
	return c.SendDocument(ctx, channelJID, docBytes, filename, mimetype, caption)
}

// SendInteractiveShopCard mengirim link etalase toko WhatsApp interaktif
func (c *Client) SendInteractiveShopCard(ctx context.Context, chat types.JID, businessOwnerJID types.JID, bodyText, buttonText string, surfaceType int32) error {
	btnParams := `{"display_text":"` + buttonText + `"}`

	msg := &waE2E.Message{
		InteractiveMessage: &waE2E.InteractiveMessage{
			Body: &waE2E.InteractiveMessage_Body{
				Text: proto.String(bodyText),
			},
			InteractiveMessage: &waE2E.InteractiveMessage_NativeFlowMessage_{
				NativeFlowMessage: &waE2E.InteractiveMessage_NativeFlowMessage{
					Buttons: []*waE2E.InteractiveMessage_NativeFlowMessage_NativeFlowButton{
						{
							Name:             proto.String("cta_shop"),
							ButtonParamsJSON: proto.String(btnParams),
						},
					},
				},
			},
		},
	}

	bizNodes := buildBizAdditionalNodes()
	_, err := c.Client.SendMessage(ctx, chat, msg, whatsmeow.SendRequestExtra{
		AdditionalNodes: &bizNodes,
	})
	return err
}
