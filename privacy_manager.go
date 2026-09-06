package meowbail

import (
	"context"
	"time"

	waBinary "go.mau.fi/whatsmeow/binary"
	"go.mau.fi/whatsmeow/types"
)

// PrivacyManager menangani konfigurasi privasi akun WhatsApp secara lengkap dan presisi
// mencakup Read Receipts (centang biru), Last Seen, Online Status, Call Add, Group Add, Profile Picture, Status, dan Disappearing Mode default.

// GetPrivacySettings mengambil konfigurasi privasi akun yang sedang aktif saat ini
func (c *Client) GetPrivacySettings(ctx context.Context) (types.PrivacySettings, error) {
	ptr, err := c.Client.TryFetchPrivacySettings(ctx, true)
	if err != nil {
		return types.PrivacySettings{}, err
	}
	if ptr == nil {
		return c.Client.GetPrivacySettings(ctx), nil
	}
	return *ptr, nil
}

// SetReadReceiptsPrivacy mengatur apakah centang biru (tanda pesan telah dibaca) diaktifkan atau dinonaktifkan
// enabled: true -> centang biru aktif ("all"), false -> centang biru mati ("none")
func (c *Client) SetReadReceiptsPrivacy(ctx context.Context, enabled bool) (types.PrivacySettings, error) {
	val := types.PrivacySettingAll
	if !enabled {
		val = types.PrivacySettingNone
	}
	return c.Client.SetPrivacySetting(ctx, types.PrivacySettingTypeReadReceipts, val)
}

// SetLastSeenPrivacy mengatur siapa yang dapat melihat status terakhir dilihat (Last Seen)
// value: types.PrivacySettingAll ("all"), types.PrivacySettingContacts ("contacts"), types.PrivacySettingNone ("none")
func (c *Client) SetLastSeenPrivacy(ctx context.Context, value types.PrivacySetting) (types.PrivacySettings, error) {
	return c.Client.SetPrivacySetting(ctx, types.PrivacySettingTypeLastSeen, value)
}

// SetOnlinePrivacy mengatur siapa yang dapat melihat saat akun sedang Online
// value: types.PrivacySettingAll ("all"), types.PrivacySettingMatchLastSeen ("match_last_seen")
func (c *Client) SetOnlinePrivacy(ctx context.Context, matchLastSeen bool) (types.PrivacySettings, error) {
	val := types.PrivacySettingAll
	if matchLastSeen {
		val = types.PrivacySettingMatchLastSeen
	}
	return c.Client.SetPrivacySetting(ctx, types.PrivacySettingTypeOnline, val)
}

// SetProfilePicturePrivacy mengatur siapa yang dapat melihat foto profil akun
func (c *Client) SetProfilePicturePrivacy(ctx context.Context, value types.PrivacySetting) (types.PrivacySettings, error) {
	return c.Client.SetPrivacySetting(ctx, types.PrivacySettingTypeProfile, value)
}

// SetStatusPrivacy mengatur siapa yang dapat melihat pembaharuan status WhatsApp
func (c *Client) SetStatusPrivacy(ctx context.Context, value types.PrivacySetting) (types.PrivacySettings, error) {
	return c.Client.SetPrivacySetting(ctx, types.PrivacySettingTypeStatus, value)
}

// SetGroupAddPrivacy mengatur siapa yang berhak menambahkan akun ke dalam grup
func (c *Client) SetGroupAddPrivacy(ctx context.Context, value types.PrivacySetting) (types.PrivacySettings, error) {
	return c.Client.SetPrivacySetting(ctx, types.PrivacySettingTypeGroupAdd, value)
}

// SetCallAddPrivacy mengatur siapa yang dapat memanggil akun via telepon/video call WhatsApp
func (c *Client) SetCallAddPrivacy(ctx context.Context, allowEveryone bool) (types.PrivacySettings, error) {
	val := types.PrivacySettingAll
	if !allowEveryone {
		val = types.PrivacySettingKnown
	}
	return c.Client.SetPrivacySetting(ctx, types.PrivacySettingTypeCallAdd, val)
}

// SetDefaultDisappearingTimer mengatur durasi pesan sementara default untuk semua obrolan baru
// Durasi umum WhatsApp: 24 jam (24 * time.Hour), 7 hari (7 * 24 * time.Hour), 90 hari (90 * 24 * time.Hour), atau 0 untuk nonaktif.
func (c *Client) SetDefaultDisappearingTimer(ctx context.Context, timer time.Duration) error {
	return c.Client.SetDefaultDisappearingTimer(ctx, timer)
}

// GetStatusPrivacy mengambil pengaturan privasi status / story broadcast WhatsApp
func (c *Client) GetStatusPrivacy(ctx context.Context) ([]types.StatusPrivacy, error) {
	return c.Client.GetStatusPrivacy(ctx)
}

// IssuePrivacyTokens mengirim permintaan untuk menerbitkan privacy token bagi target JID (trusted_contact)
func (c *Client) IssuePrivacyTokens(ctx context.Context, jids []types.JID) error {
	now := time.Now()
	for _, j := range jids {
		_, err := c.Client.DangerousInternals().IssuePrivacyToken(ctx, j, now)
		if err != nil {
			return err
		}
	}
	return nil
}

// DisappearingDurationInfo hasil kueri disappearing mode duration untuk kontak
type DisappearingDurationInfo struct {
	JID      types.JID
	Duration time.Duration
}

// UserStatusInfo hasil kueri status / bio / about pengguna WhatsApp
type UserStatusInfo struct {
	JID    types.JID
	Status string
	SetAt  time.Time
}

// FetchDisappearingDuration mengambil durasi disappearing mode aktif dari satu atau lebih kontak melalui protokol USync
// Parity dengan Baileys fetchDisappearingDuration
func (c *Client) FetchDisappearingDuration(ctx context.Context, jids []types.JID) ([]DisappearingDurationInfo, error) {
	if len(jids) == 0 {
		return nil, nil
	}

	query := []waBinary.Node{
		{Tag: "disappearing_mode"},
	}

	listNode, err := c.Client.DangerousInternals().Usync(ctx, jids, "query", "interactive", query)
	if err != nil {
		return nil, err
	}

	var results []DisappearingDurationInfo
	for _, userNode := range listNode.GetChildren() {
		if userNode.Tag != "user" {
			continue
		}
		ag := userNode.AttrGetter()
		targetJID := ag.OptionalJIDOrEmpty("jid")
		if targetJID.IsEmpty() {
			targetJID = ag.OptionalJIDOrEmpty("pn_jid")
		}

		dmNode, ok := userNode.GetOptionalChildByTag("disappearing_mode")
		var dur time.Duration
		if ok {
			durationSec := dmNode.AttrGetter().OptionalInt("duration")
			dur = time.Duration(durationSec) * time.Second
		}

		results = append(results, DisappearingDurationInfo{
			JID:      targetJID,
			Duration: dur,
		})
	}

	return results, nil
}

// UsernameInfo hasil kueri username WhatsApp (WhatsApp Usernames feature)
type UsernameInfo struct {
	JID      types.JID
	Username string
}

// FetchStatus mengambil status teks (About / Bio) dari satu atau lebih kontak melalui protokol USync
// Parity dengan Baileys fetchStatus
func (c *Client) FetchStatus(ctx context.Context, jids []types.JID) ([]UserStatusInfo, error) {
	if len(jids) == 0 {
		return nil, nil
	}

	query := []waBinary.Node{
		{Tag: "status"},
	}

	listNode, err := c.Client.DangerousInternals().Usync(ctx, jids, "query", "interactive", query)
	if err != nil {
		return nil, err
	}

	var results []UserStatusInfo
	for _, userNode := range listNode.GetChildren() {
		if userNode.Tag != "user" {
			continue
		}
		ag := userNode.AttrGetter()
		targetJID := ag.OptionalJIDOrEmpty("jid")
		if targetJID.IsEmpty() {
			targetJID = ag.OptionalJIDOrEmpty("pn_jid")
		}

		statusNode, ok := userNode.GetOptionalChildByTag("status")
		var statusText string
		var setAt time.Time
		if ok {
			contentBytes, isBytes := statusNode.Content.([]byte)
			if isBytes {
				statusText = string(contentBytes)
			}
			tSec := statusNode.AttrGetter().OptionalInt("t")
			if tSec > 0 {
				setAt = time.Unix(int64(tSec), 0)
			}
		}

		results = append(results, UserStatusInfo{
			JID:    targetJID,
			Status: statusText,
			SetAt:  setAt,
		})
	}

	return results, nil
}

// FetchUsername mengambil username publik WhatsApp dari satu atau lebih kontak melalui USync
// Parity dengan Baileys USyncUsernameProtocol
func (c *Client) FetchUsername(ctx context.Context, jids []types.JID) ([]UsernameInfo, error) {
	if len(jids) == 0 {
		return nil, nil
	}

	query := []waBinary.Node{
		{Tag: "username"},
	}

	listNode, err := c.Client.DangerousInternals().Usync(ctx, jids, "query", "interactive", query)
	if err != nil {
		return nil, err
	}

	var results []UsernameInfo
	for _, userNode := range listNode.GetChildren() {
		if userNode.Tag != "user" {
			continue
		}
		ag := userNode.AttrGetter()
		targetJID := ag.OptionalJIDOrEmpty("jid")
		if targetJID.IsEmpty() {
			targetJID = ag.OptionalJIDOrEmpty("pn_jid")
		}

		unNode, ok := userNode.GetOptionalChildByTag("username")
		var uname string
		if ok {
			contentBytes, isBytes := unNode.Content.([]byte)
			if isBytes {
				uname = string(contentBytes)
			} else if s, isStr := unNode.Content.(string); isStr {
				uname = s
			}
		}

		results = append(results, UsernameInfo{
			JID:      targetJID,
			Username: uname,
		})
	}

	return results, nil
}



