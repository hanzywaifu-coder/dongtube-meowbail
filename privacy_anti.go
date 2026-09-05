package meowbail

import (
	"context"
	"fmt"
	"time"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// ExtractQuotedMedia mengekstrak konten media dari pesan yang di-reply/quote
func (c *Client) ExtractQuotedMedia(msg *waE2E.Message) (mediaData []byte, mediaType string, err error) {
	if msg == nil {
		return nil, "", fmt.Errorf("pesan kosong")
	}

	var qm *waE2E.Message
	if msg.ExtendedTextMessage != nil && msg.ExtendedTextMessage.ContextInfo != nil {
		qm = msg.ExtendedTextMessage.ContextInfo.QuotedMessage
	}

	if qm == nil {
		return nil, "", fmt.Errorf("tidak ada pesan quoted/reply")
	}

	if qm.ImageMessage != nil {
		data, err := c.DownloadMedia(context.Background(), qm.ImageMessage)
		return data, "image", err
	}
	if qm.VideoMessage != nil {
		data, err := c.DownloadMedia(context.Background(), qm.VideoMessage)
		return data, "video", err
	}
	if qm.AudioMessage != nil {
		data, err := c.DownloadMedia(context.Background(), qm.AudioMessage)
		return data, "audio", err
	}
	if qm.StickerMessage != nil {
		data, err := c.DownloadMedia(context.Background(), qm.StickerMessage)
		return data, "sticker", err
	}
	if qm.DocumentMessage != nil {
		data, err := c.DownloadMedia(context.Background(), qm.DocumentMessage)
		return data, "document", err
	}

	return nil, "", fmt.Errorf("tipe media tidak dikenali pada quoted message")
}

// ResendViewOnce meneruskan media yang dikirim sebagai View-Once ke chat/grup tujuan dalam bentuk media biasa
func (c *Client) ResendViewOnce(ctx context.Context, targetChat types.JID, viewOnceMsg *waE2E.FutureProofMessage, caption string) error {
	if viewOnceMsg == nil || viewOnceMsg.Message == nil {
		return fmt.Errorf("pesan view-once kosong")
	}

	inner := viewOnceMsg.Message
	if inner.ImageMessage != nil {
		data, err := c.DownloadMedia(ctx, inner.ImageMessage)
		if err != nil {
			return err
		}
		return c.SendImage(ctx, targetChat, data, caption)
	}

	if inner.VideoMessage != nil {
		data, err := c.DownloadMedia(ctx, inner.VideoMessage)
		if err != nil {
			return err
		}
		return c.SendVideo(ctx, targetChat, data, caption)
	}

	if inner.AudioMessage != nil {
		data, err := c.DownloadMedia(ctx, inner.AudioMessage)
		if err != nil {
			return err
		}
		return c.SendAudio(ctx, targetChat, data)
	}

	return fmt.Errorf("tipe media view-once tidak didukung")
}

// SendGhostMention mengirim pesan teks yang men-tag seluruh pengguna secara tidak kasat mata (Ghost/Invisible Mention)
func (c *Client) SendGhostMention(ctx context.Context, chat types.JID, text string, targetJIDs []string) error {
	msg := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(text),
			ContextInfo: &waE2E.ContextInfo{
				MentionedJID: targetJIDs,
			},
		},
	}
	_, err := c.Client.SendMessage(ctx, chat, msg)
	return err
}

// SendAutoDeleteText mengirim pesan teks yang otomatis hilang (self-destructing) setelah durasi tertentu
func (c *Client) SendAutoDeleteText(ctx context.Context, chat types.JID, text string, duration time.Duration) error {
	resp, err := c.Client.SendMessage(ctx, chat, &waE2E.Message{
		Conversation: proto.String(text),
	})
	if err != nil {
		return err
	}

	go func(targetID types.MessageID) {
		time.Sleep(duration)
		_ = c.RevokeMessage(context.Background(), chat, targetID, true, types.EmptyJID)
	}(resp.ID)

	return nil
}

// FormatJID membersihkan string nomor telepon menjadi format WhatsApp JID standar
func FormatJID(phoneOrJID string) types.JID {
	if phoneOrJID == "" {
		return types.EmptyJID
	}
	clean := phoneOrJID
	for _, ch := range []string{"+", "-", " ", "@s.whatsapp.net", "@c.us", "@g.us"} {
		for {
			idx := len(clean) - len(ch)
			if idx >= 0 && clean[idx:] == ch {
				clean = clean[:idx]
			} else {
				break
			}
		}
	}
	return types.NewJID(clean, types.DefaultUserServer)
}
