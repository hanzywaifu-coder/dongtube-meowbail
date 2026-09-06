package meowbail

import (
	"context"
	"time"

	"go.mau.fi/whatsmeow/types"
)

// DisappearingTimerPresets konstanta durasi pesan sementara standar resmi WhatsApp
const (
	DisappearingTimerOff    time.Duration = 0
	DisappearingTimer24Hour time.Duration = 24 * time.Hour
	DisappearingTimer7Day   time.Duration = 7 * 24 * time.Hour
	DisappearingTimer90Day  time.Duration = 90 * 24 * time.Hour
)

// SetChatDisappearingTimer mengatur pesan sementara (disappearing messages / ephemeral)
// untuk obrolan pribadi (1-on-1 DM) ataupun obrolan grup dengan paritas protokol WhatsApp resmi
func (c *Client) SetChatDisappearingTimer(ctx context.Context, chat types.JID, timer time.Duration) error {
	return c.Client.SetDisappearingTimer(ctx, chat, timer, time.Now())
}

// DisableChatDisappearingTimer mematikan fitur pesan sementara untuk obrolan tertentu
func (c *Client) DisableChatDisappearingTimer(ctx context.Context, chat types.JID) error {
	return c.Client.SetDisappearingTimer(ctx, chat, DisappearingTimerOff, time.Now())
}
