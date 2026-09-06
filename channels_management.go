package meowbail

import (
	"context"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
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

// NewsletterGetInfoWithInvite mengambil informasi saluran dari invite code / link saluran (whatsapp.com/channel/xxx)
func (c *Client) NewsletterGetInfoWithInvite(ctx context.Context, inviteCodeOrLink string) (*types.NewsletterMetadata, error) {
	code := inviteCodeOrLink
	if strings.Contains(code, "whatsapp.com/channel/") {
		parts := strings.Split(code, "whatsapp.com/channel/")
		if len(parts) > 1 {
			code = strings.Split(parts[1], "/")[0]
			code = strings.Split(code, "?")[0]
		}
	} else if strings.Contains(code, "wa.me/channel/") {
		parts := strings.Split(code, "wa.me/channel/")
		if len(parts) > 1 {
			code = strings.Split(parts[1], "/")[0]
			code = strings.Split(code, "?")[0]
		}
	}
	code = strings.TrimSpace(code)
	return c.Client.GetNewsletterInfoWithInvite(ctx, code)
}

// NewsletterCreate membuat saluran baru dengan nama dan deskripsi
func (c *Client) NewsletterCreate(ctx context.Context, name, description string, avatarJPEG []byte) (*types.NewsletterMetadata, error) {
	return c.Client.CreateNewsletter(ctx, whatsmeow.CreateNewsletterParams{
		Name:        name,
		Description: description,
		Picture:     avatarJPEG,
	})
}

// NewsletterMarkViewed menandai postingan saluran telah dibaca/dilihat
func (c *Client) NewsletterMarkViewed(ctx context.Context, newsletterJID types.JID, serverIDs []types.MessageServerID) error {
	return c.Client.NewsletterMarkViewed(ctx, newsletterJID, serverIDs)
}

// GetSubscribedNewsletters mengambil seluruh daftar Saluran/Newsletter yang diikuti
func (c *Client) GetSubscribedNewsletters(ctx context.Context) ([]*types.NewsletterMetadata, error) {
	return c.Client.GetSubscribedNewsletters(ctx)
}

// NewsletterSubscribeLiveUpdates berlangganan live update sementara untuk suatu saluran
func (c *Client) NewsletterSubscribeLiveUpdates(ctx context.Context, jid types.JID) (time.Duration, error) {
	return c.Client.NewsletterSubscribeLiveUpdates(ctx, jid)
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
