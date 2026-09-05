package meowbail

import (
	"context"

	"go.mau.fi/whatsmeow/types"
)

// SetGroupPhoto mengubah foto/ikon grup (avatar)
func (c *Client) SetGroupPhoto(ctx context.Context, groupJID types.JID, jpegData []byte) (string, error) {
	return c.Client.SetGroupPhoto(ctx, groupJID, jpegData)
}

// RemoveGroupPhoto menghapus foto/ikon grup
func (c *Client) RemoveGroupPhoto(ctx context.Context, groupJID types.JID) error {
	_, err := c.Client.SetGroupPhoto(ctx, groupJID, nil)
	return err
}

// SetGroupJoinApprovalMode mengatur persetujuan admin untuk member baru yang mau bergabung (approve mode)
func (c *Client) SetGroupJoinApprovalMode(ctx context.Context, groupJID types.JID, requireApproval bool) error {
	return c.Client.SetGroupJoinApprovalMode(ctx, groupJID, requireApproval)
}

// SetGroupMemberAddMode mengatur siapa saja yang berhak menambahkan member ke grup (admin only atau all members)
func (c *Client) SetGroupMemberAddMode(ctx context.Context, groupJID types.JID, adminOnly bool) error {
	mode := types.GroupMemberAddModeAllMember
	if adminOnly {
		mode = types.GroupMemberAddModeAdmin
	}
	return c.Client.SetGroupMemberAddMode(ctx, groupJID, mode)
}

// SetAboutStatus mengubah status teks bio profil akun WhatsApp ("About")
func (c *Client) SetAboutStatus(ctx context.Context, statusText string) error {
	return c.Client.SetStatusMessage(ctx, types.SetStatusInput{
		Text: &statusText,
	})
}