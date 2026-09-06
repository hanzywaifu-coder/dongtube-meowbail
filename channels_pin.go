package meowbail

import (
	"context"
	"time"

	"go.mau.fi/whatsmeow/appstate"
	"go.mau.fi/whatsmeow/proto/waCommon"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// PinAction jenis aksi pin
type PinAction int32

const (
	PinActionPin   PinAction = 1
	PinActionUnpin PinAction = 2
)

// PinDuration durasi pesan di-pin (standar WhatsApp: 24h, 7d, 30d)
type PinDuration uint32

const (
	PinDuration24Hours PinDuration = 86400
	PinDuration7Days   PinDuration = 604800
	PinDuration30Days  PinDuration = 2592000
)

// PinMessage menyematkan atau melepas sematan pesan dalam percakapan
func (c *Client) PinMessage(ctx context.Context, chat types.JID, targetMsgID types.MessageID, action PinAction, duration PinDuration) error {
	var pinType waE2E.PinInChatMessage_Type
	if action == PinActionPin {
		pinType = waE2E.PinInChatMessage_PIN_FOR_ALL
	} else {
		pinType = waE2E.PinInChatMessage_UNPIN_FOR_ALL
	}

	if duration == 0 {
		duration = PinDuration7Days
	}

	msg := &waE2E.Message{
		PinInChatMessage: &waE2E.PinInChatMessage{
			Key: &waCommon.MessageKey{
				RemoteJID: proto.String(chat.String()),
				ID:        proto.String(string(targetMsgID)),
			},
			Type:              pinType.Enum(),
			SenderTimestampMS: proto.Int64(0),
		},
	}

	_, err := c.Client.SendMessage(ctx, chat, msg)
	return err
}

// PinChat menyematkan atau melepas sematan obrolan di daftar chat utama (App State Sync)
func (c *Client) PinChat(ctx context.Context, chat types.JID, pin bool) error {
	patch := appstate.BuildPin(chat, pin)
	return c.Client.SendAppState(ctx, patch)
}

// MuteChat membisukan atau mengaktifkan kembali notifikasi obrolan untuk durasi tertentu (0 = selamanya)
func (c *Client) MuteChat(ctx context.Context, chat types.JID, mute bool, duration time.Duration) error {
	patch := appstate.BuildMute(chat, mute, duration)
	return c.Client.SendAppState(ctx, patch)
}

// ArchiveChat mengarsipkan atau membuka arsip obrolan
func (c *Client) ArchiveChat(ctx context.Context, chat types.JID, archive bool) error {
	patch := appstate.BuildArchive(chat, archive, time.Now(), nil)
	return c.Client.SendAppState(ctx, patch)
}

// MarkChatAsRead menandai obrolan sebagai sudah dibaca atau belum dibaca
func (c *Client) MarkChatAsRead(ctx context.Context, chat types.JID, read bool) error {
	patch := appstate.BuildMarkChatAsRead(chat, read, time.Now(), nil)
	return c.Client.SendAppState(ctx, patch)
}

// NewsletterFollow mengikuti newsletter/saluran WhatsApp
func (c *Client) NewsletterFollow(ctx context.Context, channelJID types.JID) error {
	return c.Client.FollowNewsletter(ctx, channelJID)
}

// NewsletterUnfollow berhenti mengikuti saluran
func (c *Client) NewsletterUnfollow(ctx context.Context, channelJID types.JID) error {
	return c.Client.UnfollowNewsletter(ctx, channelJID)
}

// NewsletterMute membisukan atau membunyikan notifikasi saluran
func (c *Client) NewsletterMute(ctx context.Context, channelJID types.JID, mute bool) error {
	if mute {
		return c.Client.NewsletterToggleMute(ctx, channelJID, true)
	}
	return c.Client.NewsletterToggleMute(ctx, channelJID, false)
}

// MarkChatAsUnread menandai obrolan sebagai belum dibaca (Unread)
func (c *Client) MarkChatAsUnread(ctx context.Context, chat types.JID) error {
	return c.MarkChatAsRead(ctx, chat, false)
}
