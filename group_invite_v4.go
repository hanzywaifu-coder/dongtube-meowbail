package meowbail

import (
	"context"
	"fmt"
	"strconv"

	waBinary "go.mau.fi/whatsmeow/binary"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
)

// GroupAcceptInviteV4 menerima undangan masuk grup v4 (GroupInviteMessage)
// Mengikuti implementasi Baileys groupAcceptInviteV4
func (c *Client) GroupAcceptInviteV4(ctx context.Context, groupJID types.JID, inviter types.JID, inviteCode string, inviteExpiration int64) (types.JID, error) {
	if inviteCode == "" {
		return types.EmptyJID, fmt.Errorf("invite code tidak boleh kosong")
	}

	attrs := waBinary.Attrs{
		"code":  inviteCode,
		"admin": inviter.String(),
	}
	if inviteExpiration > 0 {
		attrs["expiration"] = strconv.FormatInt(inviteExpiration, 10)
	}

	queryNode := waBinary.Node{
		Tag:   "accept",
		Attrs: attrs,
	}

	resp, err := c.Client.DangerousInternals().SendGroupIQ(ctx, "set", groupJID, queryNode)
	if err != nil {
		return types.EmptyJID, err
	}

	if resp != nil {
		return groupJID, nil
	}

	return groupJID, nil
}

// GroupAcceptInviteV4FromMessage mem-parse langsung waE2E.GroupInviteMessage dan mengeksekusi accept invite
func (c *Client) GroupAcceptInviteV4FromMessage(ctx context.Context, inviter types.JID, msg *waE2E.GroupInviteMessage) (types.JID, error) {
	if msg == nil {
		return types.EmptyJID, fmt.Errorf("pesan group invite kosong")
	}

	groupJID, err := types.ParseJID(msg.GetGroupJID())
	if err != nil {
		return types.EmptyJID, fmt.Errorf("invalid group JID: %w", err)
	}

	return c.GroupAcceptInviteV4(ctx, groupJID, inviter, msg.GetInviteCode(), msg.GetInviteExpiration())
}
