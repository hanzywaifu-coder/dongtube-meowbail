package meowbail

import (
	"context"

	"go.mau.fi/whatsmeow"
	waBinary "go.mau.fi/whatsmeow/binary"
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

// GetBlocklist mengambil seluruh daftar kontak yang sedang diblokir
func (c *Client) GetBlocklist(ctx context.Context) (*types.Blocklist, error) {
	return c.Client.GetBlocklist(ctx)
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

// CreateCommunity membuat Komunitas (Parent Group) baru di WhatsApp
func (c *Client) CreateCommunity(ctx context.Context, name, description string) (*types.GroupInfo, error) {
	req := whatsmeow.ReqCreateGroup{
		Name: name,
		GroupParent: types.GroupParent{
			IsParent: true,
		},
	}
	info, err := c.Client.CreateGroup(ctx, req)
	if err != nil {
		return nil, err
	}

	if description != "" && info != nil {
		_ = c.Client.SetGroupTopic(ctx, info.JID, "", "", description)
	}

	return info, nil
}

// LinkSubgroupToCommunity menghubungkan grup biasa ke dalam komunitas
func (c *Client) LinkSubgroupToCommunity(ctx context.Context, communityJID types.JID, groupJID types.JID) error {
	return c.Client.LinkGroup(ctx, communityJID, groupJID)
}

// UnlinkSubgroupFromCommunity memutuskan grup dari komunitas
func (c *Client) UnlinkSubgroupFromCommunity(ctx context.Context, communityJID types.JID, groupJID types.JID) error {
	return c.Client.UnlinkGroup(ctx, communityJID, groupJID)
}

// GetCommunitySubGroups mengambil seluruh grup anak / subgroup yang terhubung ke komunitas
func (c *Client) GetCommunitySubGroups(ctx context.Context, communityJID types.JID) ([]*types.GroupLinkTarget, error) {
	return c.Client.GetSubGroups(ctx, communityJID)
}

// GetCommunityParticipants mengambil seluruh partisipan dari semua grup dalam komunitas
func (c *Client) GetCommunityParticipants(ctx context.Context, communityJID types.JID) ([]types.JID, error) {
	return c.Client.GetLinkedGroupsParticipants(ctx, communityJID)
}

// CommunityUpdateSubject mengubah nama / judul Komunitas WhatsApp
func (c *Client) CommunityUpdateSubject(ctx context.Context, communityJID types.JID, newSubject string) error {
	return c.Client.SetGroupName(ctx, communityJID, newSubject)
}

// CommunityUpdateDescription mengubah deskripsi / topik Komunitas WhatsApp
func (c *Client) CommunityUpdateDescription(ctx context.Context, communityJID types.JID, newDescription string) error {
	return c.Client.SetGroupTopic(ctx, communityJID, "", "", newDescription)
}

// CommunityLeave keluar dari Komunitas WhatsApp
func (c *Client) CommunityLeave(ctx context.Context, communityJID types.JID) error {
	return c.Client.LeaveGroup(ctx, communityJID)
}

// CommunityGetInviteLink mengambil kode link undangan Komunitas WhatsApp
func (c *Client) CommunityGetInviteLink(ctx context.Context, communityJID types.JID, reset bool) (string, error) {
	return c.Client.GetGroupInviteLink(ctx, communityJID, reset)
}

// CommunityAcceptInvite bergabung ke Komunitas WhatsApp via invite code
func (c *Client) CommunityAcceptInvite(ctx context.Context, code string) (types.JID, error) {
	return c.Client.JoinGroupWithLink(ctx, code)
}

// CommunityParticipantsUpdate memperbarui status partisipan di komunitas (promote, demote, remove dari seluruh grup terkait)
func (c *Client) CommunityParticipantsUpdate(ctx context.Context, communityJID types.JID, participants []types.JID, action whatsmeow.ParticipantChange) ([]types.GroupParticipant, error) {
	if action == whatsmeow.ParticipantChangeRemove {
		// Remove dari komunitas sekaligus linked groups
		content := make([]waBinary.Node, len(participants))
		for i, p := range participants {
			content[i] = waBinary.Node{
				Tag:   "participant",
				Attrs: waBinary.Attrs{"jid": p},
			}
		}
		resp, err := c.Client.DangerousInternals().SendGroupIQ(ctx, "set", communityJID, waBinary.Node{
			Tag:     "remove",
			Attrs:   waBinary.Attrs{"linked_groups": "true"},
			Content: content,
		})
		if err != nil {
			return nil, err
		}
		if resp != nil {
			removeNode, ok := resp.GetOptionalChildByTag("remove")
			if ok {
				participantNodes := removeNode.GetChildrenByTag("participant")
				res := make([]types.GroupParticipant, len(participantNodes))
				for i, node := range participantNodes {
					res[i] = types.GroupParticipant{
						JID: node.AttrGetter().JID("jid"),
					}
				}
				return res, nil
			}
		}
		return nil, nil
	}

	return c.Client.UpdateGroupParticipants(ctx, communityJID, participants, action)
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
