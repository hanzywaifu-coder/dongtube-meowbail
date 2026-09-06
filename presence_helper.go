package meowbail

import (
	"context"
	"time"

	"go.mau.fi/whatsmeow/types"
)

// SimulateTyping menyimulasikan bot sedang mengetik teks selama durasi tertentu
func (c *Client) SimulateTyping(ctx context.Context, chat types.JID, duration time.Duration) error {
	_ = c.Client.SendChatPresence(ctx, chat, types.ChatPresenceComposing, types.ChatPresenceMediaText)
	select {
	case <-time.After(duration):
	case <-ctx.Done():
		_ = c.Client.SendChatPresence(context.Background(), chat, types.ChatPresencePaused, types.ChatPresenceMediaText)
		return ctx.Err()
	}
	return c.Client.SendChatPresence(ctx, chat, types.ChatPresencePaused, types.ChatPresenceMediaText)
}

// SimulateRecording menyimulasikan bot sedang merekam audio/voice note selama durasi tertentu
func (c *Client) SimulateRecording(ctx context.Context, chat types.JID, duration time.Duration) error {
	_ = c.Client.SendChatPresence(ctx, chat, types.ChatPresenceComposing, types.ChatPresenceMediaAudio)
	select {
	case <-time.After(duration):
	case <-ctx.Done():
		_ = c.Client.SendChatPresence(context.Background(), chat, types.ChatPresencePaused, types.ChatPresenceMediaAudio)
		return ctx.Err()
	}
	return c.Client.SendChatPresence(ctx, chat, types.ChatPresencePaused, types.ChatPresenceMediaAudio)
}

// MarkReadSimple menandai pesan-pesan terakhir dalam obrolan sebagai sudah dibaca (Read Receipts)
func (c *Client) MarkReadSimple(ctx context.Context, chat types.JID, messageIDs []types.MessageID) error {
	return c.Client.MarkRead(ctx, messageIDs, time.Now(), chat, types.EmptyJID)
}
