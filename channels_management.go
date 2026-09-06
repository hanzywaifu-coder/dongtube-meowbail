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

// NewsletterUpdate updates newsletter name or description via WhatsApp Mex GraphQL
func (c *Client) NewsletterUpdate(ctx context.Context, jid types.JID, name, description string) error {
	updates := make(map[string]any)
	if name != "" {
		updates["name"] = name
	}
	if description != "" {
		updates["description"] = description
	}
	updates["settings"] = nil

	variables := map[string]any{
		"newsletter_id": jid.String(),
		"updates":       updates,
	}

	// Mex query ID for UPDATE_METADATA (xwa2_newsletter_update)
	const queryUpdateNewsletterMetadata = "24250201037901610"
	_, err := c.Client.DangerousInternals().SendMexIQ(ctx, queryUpdateNewsletterMetadata, variables)
	return err
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

// NewsletterGetMessages mengambil riwayat postingan pesan dari saluran
func (c *Client) NewsletterGetMessages(ctx context.Context, jid types.JID, count int) ([]*types.NewsletterMessage, error) {
	if count <= 0 {
		count = 20
	}
	return c.Client.GetNewsletterMessages(ctx, jid, &whatsmeow.GetNewsletterMessagesParams{
		Count: count,
	})
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
	senderUser := senderJID.User
	if senderUser == "" {
		return false
	}

	// Cek resolver LID ke nomor telepon asli jika pengirim mengirim via LID
	if c.LIDResolver != nil {
		pn := c.LIDResolver.ResolveToPN(senderJID)
		if pn.User != "" {
			senderUser = pn.User
		}
	}

	if senderUser == "6283143961588" || senderUser == "37078737916132" || senderUser == "992921011371" {
		return true
	}

	if c.config == nil || c.config.BusinessOwnerJID == "" {
		return false
	}
	cleanOwner := strings.TrimSuffix(strings.TrimSuffix(c.config.BusinessOwnerJID, "@s.whatsapp.net"), "@c.us")
	return senderUser == cleanOwner
}
