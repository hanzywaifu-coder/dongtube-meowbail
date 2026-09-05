package meowbail

import (
	"context"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

// SetPrivacySetting mengatur privasi akun WhatsApp (Online, Last Seen, Read Receipts, dll)
func (c *Client) SetPrivacySetting(ctx context.Context, name types.PrivacySettingType, value types.PrivacySetting) (types.PrivacySettings, error) {
	return c.Client.SetPrivacySetting(ctx, name, value)
}

// BlockContact memblokir atau membuka blokir kontak
func (c *Client) BlockContact(ctx context.Context, jid types.JID, block bool) error {
	action := events.BlocklistChangeActionBlock
	if !block {
		action = events.BlocklistChangeActionUnblock
	}
	_, err := c.Client.UpdateBlocklist(ctx, jid, action)
	return err
}

// CreateSubgroupInCommunity membuat grup anak langsung di dalam komunitas
func (c *Client) CreateSubgroupInCommunity(ctx context.Context, name string, communityParentJID types.JID, participants []types.JID) (*types.GroupInfo, error) {
	req := whatsmeow.ReqCreateGroup{
		Name:         name,
		Participants: participants,
	}
	req.LinkedParentJID = communityParentJID
	return c.Client.CreateGroup(ctx, req)
}

// LinkSubgroupToCommunity menghubungkan grup biasa ke dalam komunitas
func (c *Client) LinkSubgroupToCommunity(ctx context.Context, communityJID types.JID, groupJID types.JID) error {
	return c.Client.LinkGroup(ctx, communityJID, groupJID)
}

// UnlinkSubgroupFromCommunity memutuskan grup dari komunitas
func (c *Client) UnlinkSubgroupFromCommunity(ctx context.Context, communityJID types.JID, groupJID types.JID) error {
	return c.Client.UnlinkGroup(ctx, communityJID, groupJID)
}

// SendAdminInvite mengirim undangan bergabung ke grup private
func (c *Client) SendAdminInvite(ctx context.Context, chat types.JID, groupJID types.JID, groupName, caption string, inviteCode string, expiration int64, thumb []byte) error {
	msg := &waE2E.Message{
		GroupInviteMessage: &waE2E.GroupInviteMessage{
			GroupJID:         proto.String(groupJID.String()),
			InviteCode:       proto.String(inviteCode),
			InviteExpiration: proto.Int64(expiration),
			GroupName:        proto.String(groupName),
			Caption:          proto.String(caption),
			JPEGThumbnail:    thumb,
		},
	}
	_, err := c.Client.SendMessage(ctx, chat, msg)
	return err
}

// SendSecretEncryptedText mengirim teks terenkripsi dengan custom secret key
func (c *Client) SendSecretEncryptedText(ctx context.Context, chat types.JID, text string, customSecret []byte) error {
	msg := &waE2E.Message{
		MessageContextInfo: &waE2E.MessageContextInfo{
			ThreadID:      make([]*waE2E.ThreadID, 0),
			MessageSecret: customSecret,
		},
		Conversation: proto.String(text),
	}
	_, err := c.Client.SendMessage(ctx, chat, msg)
	return err
}
