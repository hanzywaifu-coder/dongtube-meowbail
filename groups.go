package meowbail

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

// GetGroupInfo gets group information
func (c *Client) GetGroupInfo(ctx context.Context, groupJID types.JID) (*types.GroupInfo, error) {
	return c.Client.GetGroupInfo(ctx, groupJID)
}

// GetJoinedGroups mengambil seluruh daftar grup tempat bot bergabung beserta metadata lengkapnya
func (c *Client) GetJoinedGroups(ctx context.Context) ([]*types.GroupInfo, error) {
	return c.Client.GetJoinedGroups(ctx)
}

// GetGroupInviteLink gets group invite link
func (c *Client) GetGroupInviteLink(ctx context.Context, groupJID types.JID, reset bool) (string, error) {
	return c.Client.GetGroupInviteLink(ctx, groupJID, reset)
}

// GetGroupInfoFromInviteLink mengambil metadata grup dari link undangan chat.whatsapp.com/xxx
func (c *Client) GetGroupInfoFromInviteLink(ctx context.Context, linkOrCode string) (*types.GroupInfo, error) {
	code := linkOrCode
	if idx := strings.Index(code, "chat.whatsapp.com/"); idx != -1 {
		code = code[idx+len("chat.whatsapp.com/"):]
		if slashIdx := strings.Index(code, "/"); slashIdx != -1 {
			code = code[:slashIdx]
		}
		if qIdx := strings.Index(code, "?"); qIdx != -1 {
			code = code[:qIdx]
		}
	}
	code = strings.TrimSpace(code)
	return c.Client.GetGroupInfoFromLink(ctx, code)
}

// JoinGroupWithInviteLink bergabung ke grup menggunakan link undangan
func (c *Client) JoinGroupWithInviteLink(ctx context.Context, linkOrCode string) (types.JID, error) {
	code := linkOrCode
	if idx := strings.Index(code, "chat.whatsapp.com/"); idx != -1 {
		code = code[idx+len("chat.whatsapp.com/"):]
		if slashIdx := strings.Index(code, "/"); slashIdx != -1 {
			code = code[:slashIdx]
		}
		if qIdx := strings.Index(code, "?"); qIdx != -1 {
			code = code[:qIdx]
		}
	}
	code = strings.TrimSpace(code)
	return c.Client.JoinGroupWithLink(ctx, code)
}

// SetGroupName sets group name
func (c *Client) SetGroupName(ctx context.Context, groupJID types.JID, name string) error {
	return c.Client.SetGroupName(ctx, groupJID, name)
}

// SetGroupDescription sets group description
func (c *Client) SetGroupDescription(ctx context.Context, groupJID types.JID, description string) error {
	return c.Client.SetGroupDescription(ctx, groupJID, description)
}

// SetGroupAnnounce sets group announce mode
func (c *Client) SetGroupAnnounce(ctx context.Context, groupJID types.JID, announce bool) error {
	return c.Client.SetGroupAnnounce(ctx, groupJID, announce)
}

// SetGroupLocked sets group locked status
func (c *Client) SetGroupLocked(ctx context.Context, groupJID types.JID, locked bool) error {
	return c.Client.SetGroupLocked(ctx, groupJID, locked)
}

// SetGroupEphemeral sets group disappearing messages timer
func (c *Client) SetGroupEphemeral(ctx context.Context, groupJID types.JID, timer time.Duration) error {
	return c.Client.SetDisappearingTimer(ctx, groupJID, timer, time.Now())
}

// PromoteParticipant promotes a participant to admin
func (c *Client) PromoteParticipant(ctx context.Context, groupJID types.JID, participants ...types.JID) error {
	_, err := c.Client.UpdateGroupParticipants(ctx, groupJID, participants, whatsmeow.ParticipantChangePromote)
	return err
}

// DemoteParticipant demotes an admin to regular member
func (c *Client) DemoteParticipant(ctx context.Context, groupJID types.JID, participants ...types.JID) error {
	_, err := c.Client.UpdateGroupParticipants(ctx, groupJID, participants, whatsmeow.ParticipantChangeDemote)
	return err
}

// RemoveParticipant removes a participant from group
func (c *Client) RemoveParticipant(ctx context.Context, groupJID types.JID, participants ...types.JID) error {
	_, err := c.Client.UpdateGroupParticipants(ctx, groupJID, participants, whatsmeow.ParticipantChangeRemove)
	return err
}

// AddParticipant adds a participant to group
func (c *Client) AddParticipant(ctx context.Context, groupJID types.JID, participants ...types.JID) error {
	_, err := c.Client.UpdateGroupParticipants(ctx, groupJID, participants, whatsmeow.ParticipantChangeAdd)
	return err
}

// LeaveGroup makes the bot leave the group
func (c *Client) LeaveGroup(ctx context.Context, groupJID types.JID) error {
	return c.Client.LeaveGroup(ctx, groupJID)
}

// IsBotAdmin checks if the bot is admin in the group
func (c *Client) IsBotAdmin(ctx context.Context, groupJID types.JID) (bool, error) {
	groupInfo, err := c.Client.GetGroupInfo(ctx, groupJID)
	if err != nil {
		return false, err
	}

	// Check both ID (phone) and LID (linked device)
	botJID := c.Client.Store.ID.User
	botLID := c.Client.Store.LID.User
	for _, p := range groupInfo.Participants {
		if (p.JID.User == botJID || p.JID.User == botLID) && p.IsAdmin {
			return true, nil
		}
	}

	return false, nil
}

// IsUserAdmin checks if a user is admin in the group
func (c *Client) IsUserAdmin(ctx context.Context, groupJID types.JID, userJID types.JID) (bool, error) {
	groupInfo, err := c.Client.GetGroupInfo(ctx, groupJID)
	if err != nil {
		return false, err
	}

	for _, p := range groupInfo.Participants {
		if p.JID == userJID && p.IsAdmin {
			return true, nil
		}
	}

	return false, nil
}

// GetGroupParticipants returns all participants in the group
func (c *Client) GetGroupParticipants(ctx context.Context, groupJID types.JID) ([]types.GroupParticipant, error) {
	groupInfo, err := c.Client.GetGroupInfo(ctx, groupJID)
	if err != nil {
		return nil, err
	}

	return groupInfo.Participants, nil
}

// TagAll mentions all participants in the group
func (c *Client) TagAll(ctx context.Context, groupJID types.JID, message string) error {
	participants, err := c.GetGroupParticipants(ctx, groupJID)
	if err != nil {
		return err
	}

	var mentionedJids []string
	for _, p := range participants {
		mentionedJids = append(mentionedJids, p.JID.String())
	}

	text := message + "\n"
	for _, p := range participants {
		text += fmt.Sprintf("@%s\n", p.JID.User)
	}

	return c.SendText(ctx, groupJID, text, &MessageOptions{
		Mentions: mentionedJids,
	})
}
