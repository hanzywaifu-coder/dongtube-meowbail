package meowbail

import (
	"context"

	waBinary "go.mau.fi/whatsmeow/binary"
	"go.mau.fi/whatsmeow/types"
)

// CleanDirtyBits membersihkan bit sinkronisasi kotor (dirty bits) untuk sinkronisasi akun atau grup
// Mengikuti Baileys cleanDirtyBits (iq set xmlns="urn:xmpp:whatsapp:dirty")
func (c *Client) CleanDirtyBits(ctx context.Context, dirtyType string) error {
	cleanNode := waBinary.Node{
		Tag: "clean",
		Attrs: waBinary.Attrs{
			"type": dirtyType,
		},
	}

	queryNode := waBinary.Node{
		Tag: "iq",
		Attrs: waBinary.Attrs{
			"to":    types.ServerJID.String(),
			"type":  "set",
			"xmlns": "urn:xmpp:whatsapp:dirty",
		},
		Content: []waBinary.Node{cleanNode},
	}

	return c.Client.DangerousInternals().SendNode(ctx, queryNode)
}
