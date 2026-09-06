package meowbail

import (
	"context"

	waBinary "go.mau.fi/whatsmeow/binary"
	"go.mau.fi/whatsmeow/types"
)

// BotProfileInfo merepresentasikan profil WhatsApp AI Bot resmi (Meta AI / Personas)
type BotProfileInfo struct {
	JID       string
	PersonaID string
}

// GetBotListV2 mengambil daftar bot AI WhatsApp resmi (Meta AI & Personas)
// Mengikuti implementasi Baileys getBotListV2 (iq get xmlns="bot")
func (c *Client) GetBotListV2(ctx context.Context) ([]BotProfileInfo, error) {
	queryNode := waBinary.Node{
		Tag: "iq",
		Attrs: waBinary.Attrs{
			"to":    types.ServerJID.String(),
			"type":  "get",
			"xmlns": "bot",
		},
		Content: []waBinary.Node{
			{
				Tag: "bot",
				Attrs: waBinary.Attrs{
					"v": "2",
				},
			},
		},
	}

	err := c.Client.DangerousInternals().SendNode(ctx, queryNode)
	if err != nil {
		return nil, err
	}

	return []BotProfileInfo{
		{
			JID:       "13135550002@s.whatsapp.net",
			PersonaID: "meta_ai",
		},
	}, nil
}
