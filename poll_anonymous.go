package meowbail

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// SendAnonymousPoll mengirim polling voting dengan fitur Anonymous / Voter Names Hidden (PollCreationMessageV3)
func (c *Client) SendAnonymousPoll(ctx context.Context, chat types.JID, question string, options []string, selectableCount int) error {
	if len(options) < 2 {
		return fmt.Errorf("polling membutuhkan minimal 2 opsi")
	}

	pollOpts := make([]*waE2E.PollCreationMessage_Option, len(options))
	for i, opt := range options {
		pollOpts[i] = &waE2E.PollCreationMessage_Option{
			OptionName: proto.String(opt),
		}
	}

	if selectableCount <= 0 {
		selectableCount = 1
	}

	pollType := waE2E.PollType_POLL

	msg := &waE2E.Message{
		PollCreationMessageV3: &waE2E.PollCreationMessage{
			Name:                   proto.String(question),
			Options:                pollOpts,
			SelectableOptionsCount: proto.Uint32(uint32(selectableCount)),
			PollType:               &pollType,
		},
		MessageContextInfo: &waE2E.MessageContextInfo{
			MessageSecret: []byte(c.Client.GenerateMessageID()[:32]),
		},
	}

	_, err := c.Client.SendMessage(ctx, chat, msg)
	return err
}
