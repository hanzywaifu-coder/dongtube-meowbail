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

// GetUserProfiles mengambil profil lengkap satu atau beberapa nomor (status bio, avatar ID, verified business name, daftar perangkat)
func (c *Client) GetUserProfiles(ctx context.Context, jids []types.JID) (map[types.JID]types.UserInfo, error) {
	return c.Client.GetUserInfo(ctx, jids)
}

// GetUserProfile mengambil profil tunggal seorang pengguna
func (c *Client) GetUserProfile(ctx context.Context, jid types.JID) (*types.UserInfo, error) {
	res, err := c.Client.GetUserInfo(ctx, []types.JID{jid})
	if err != nil {
		return nil, err
	}
	if info, ok := res[jid]; ok {
		return &info, nil
	}
	return nil, nil
}