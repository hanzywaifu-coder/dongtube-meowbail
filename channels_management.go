package meowbail

import (
	"context"
	"strings"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// NewsletterPost mengirim postingan baru ke Saluran / Newsletter milik bot/admin
func (c *Client) NewsletterPost(ctx context.Context, newsletterJID types.JID, text string) (*waE2E.Message, error) {
	msg := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(text),
		},
	}
	_, err := c.Client.SendMessage(ctx, newsletterJID, msg)
	return msg, err
}

// NewsletterPostImage mengirim gambar ke Saluran / Newsletter
func (c *Client) NewsletterPostImage(ctx context.Context, newsletterJID types.JID, imageBytes []byte, caption string) error {
	return c.SendImage(ctx, newsletterJID, imageBytes, caption)
}

// NewsletterGetMetadata mengambil informasi dan deskripsi saluran
func (c *Client) NewsletterGetMetadata(ctx context.Context, newsletterJID types.JID) (*types.NewsletterMetadata, error) {
	return c.Client.GetNewsletterInfo(ctx, newsletterJID)
}

// NewsletterMarkViewed menandai postingan saluran telah dibaca/dilihat
func (c *Client) NewsletterMarkViewed(ctx context.Context, newsletterJID types.JID, serverIDs []types.MessageServerID) error {
	return c.Client.NewsletterMarkViewed(ctx, newsletterJID, serverIDs)
}

// FormatMention membersihkan teks dan menghasilkan tag string @user dan slice MentionedJID
func FormatMention(textWithAt string, participants []types.GroupParticipant) (cleanText string, mentionedJIDs []string) {
	for _, p := range participants {
		user := p.JID.User
		if strings.Contains(textWithAt, "@"+user) {
			mentionedJIDs = append(mentionedJIDs, p.JID.ToNonAD().String())
		}
	}
	return textWithAt, mentionedJIDs
}

// CheckBotAdmin mengecek apakah nomor bot sendiri adalah admin di grup tujuan
func (c *Client) CheckBotAdmin(ctx context.Context, groupJID types.JID) (bool, error) {
	return c.IsBotAdmin(ctx, groupJID)
}

// IsSenderOwner mengecek apakah pengirim adalah pemilik bot
func (c *Client) IsSenderOwner(senderJID types.JID) bool {
	if c.config == nil || c.config.BusinessOwnerJID == "" {
		return false
	}
	cleanOwner := strings.TrimSuffix(strings.TrimSuffix(c.config.BusinessOwnerJID, "@s.whatsapp.net"), "@c.us")
	return senderJID.User == cleanOwner
}
