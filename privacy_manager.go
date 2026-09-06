package meowbail

import (
	"context"
	"time"

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
