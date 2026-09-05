package meowbail

import (
	"context"

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
