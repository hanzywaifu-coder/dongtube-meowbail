package meowbail

import (
	"context"

	"go.mau.fi/whatsmeow/types"
)

// GetContactInfo mengambil informasi nama dan detail kontak dari local store whatsmeow
func (c *Client) GetContactInfo(ctx context.Context, jid types.JID) (*types.ContactInfo, error) {
	if c.Client == nil || c.Client.Store == nil || c.Client.Store.Contacts == nil {
		return &types.ContactInfo{Found: false}, nil
	}
	info, err := c.Client.Store.Contacts.GetContact(ctx, jid)
	return &info, err
}

// GetAllContacts mengambil seluruh daftar kontak yang tersimpan di database lokal
func (c *Client) GetAllContacts(ctx context.Context) (map[types.JID]types.ContactInfo, error) {
	if c.Client == nil || c.Client.Store == nil || c.Client.Store.Contacts == nil {
		return make(map[types.JID]types.ContactInfo), nil
	}
	return c.Client.Store.Contacts.GetAllContacts(ctx)
}

// GetUserDevices mengambil daftar perangkat terhubung (multi-device) untuk kumpulan kontak WhatsApp
// Parity dengan Baileys getUSyncDevices / fbid:devices
func (c *Client) GetUserDevices(ctx context.Context, jids []types.JID) ([]types.JID, error) {
	if len(jids) == 0 {
		return nil, nil
	}
	return c.Client.DangerousInternals().GetFBIDDevices(ctx, jids)
}
