package meowbail

import (
	"context"
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// ExtractText mengekstrak pesan teks dari berbagai varian payload WhatsApp
func ExtractText(msg *waE2E.Message) string {
	if msg == nil {
		return ""
	}
	if msg.Conversation != nil {
		return *msg.Conversation
	}
	if msg.ExtendedTextMessage != nil && msg.ExtendedTextMessage.Text != nil {
		return *msg.ExtendedTextMessage.Text
	}
	if msg.ImageMessage != nil && msg.ImageMessage.Caption != nil {
		return *msg.ImageMessage.Caption
	}
	if msg.VideoMessage != nil && msg.VideoMessage.Caption != nil {
		return *msg.VideoMessage.Caption
	}
	if msg.DocumentMessage != nil && msg.DocumentMessage.Caption != nil {
		return *msg.DocumentMessage.Caption
	}
	if msg.ButtonsResponseMessage != nil && msg.ButtonsResponseMessage.SelectedButtonID != nil {
		return *msg.ButtonsResponseMessage.SelectedButtonID
	}
	if msg.ListResponseMessage != nil && msg.ListResponseMessage.SingleSelectReply != nil && msg.ListResponseMessage.SingleSelectReply.SelectedRowID != nil {
		return *msg.ListResponseMessage.SingleSelectReply.SelectedRowID
	}
	return ""
}

// ExtractSender mengekstrak JID pengirim secara aman dari grup maupun DM pribadi
func ExtractSender(evt *MessageEvent) types.JID {
	if evt == nil {
		return types.EmptyJID
	}
	if evt.IsGroup {
		return evt.Sender.ToNonAD()
	}
	return evt.Chat.ToNonAD()
}

// SendSilentReply mengirim balasan tanpa suara notifikasi pada obrolan (Silent Message Notification)
func (c *Client) SendSilentReply(ctx context.Context, chat types.JID, text string) error {
	msg := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(text),
		},
	}
	_, err := c.Client.SendMessage(ctx, chat, msg)
	return err
}

// SendAutoReaction otomatis memberi reaksi emoji ke pesan masuk
func (c *Client) SendAutoReaction(ctx context.Context, chat types.JID, targetID types.MessageID, emoji string) error {
	return c.SendReaction(ctx, chat, targetID, emoji)
}

// CleanNumber membersihkan format nomor internasional menjadi digit bersih
func CleanNumber(phone string) string {
	var sb strings.Builder
	for _, ch := range phone {
		if ch >= '0' && ch <= '9' {
			sb.WriteRune(ch)
		}
	}
	return sb.String()
}

// FormatMentionList membentuk format string @user1 @user2 dan slice string jid untuk tag massal
func FormatMentionList(jids []types.JID) (text string, mentions []string) {
	var tags []string
	for _, j := range jids {
		user := j.User
		if user != "" {
			tags = append(tags, fmt.Sprintf("@%s", user))
			mentions = append(mentions, j.ToNonAD().String())
		}
	}
	return strings.Join(tags, " "), mentions
}
