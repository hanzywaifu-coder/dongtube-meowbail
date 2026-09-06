package meowbail

import (
	"context"

	waBinary "go.mau.fi/whatsmeow/binary"
	"go.mau.fi/whatsmeow/types"
)

// BusinessProfileProps mendefinisikan field profil bisnis WhatsApp
type BusinessProfileProps struct {
	Address     string
	Email       string
	Description string
	Websites    []string
}

// UpdateBusinessProfile memperbarui profil bisnis akun WhatsApp (alamat, email, deskripsi, website)
// Mengikuti implementasi Baileys updateBussinesProfile (iq set xmlns="w:biz")
func (c *Client) UpdateBusinessProfile(ctx context.Context, props BusinessProfileProps) error {
	var nodes []waBinary.Node

	if props.Address != "" {
		nodes = append(nodes, waBinary.Node{
			Tag:     "address",
			Content: []byte(props.Address),
		})
	}
	if props.Email != "" {
		nodes = append(nodes, waBinary.Node{
			Tag:     "email",
			Content: []byte(props.Email),
		})
	}
	if props.Description != "" {
		nodes = append(nodes, waBinary.Node{
			Tag:     "description",
			Content: []byte(props.Description),
		})
	}
	for _, web := range props.Websites {
		if web != "" {
			nodes = append(nodes, waBinary.Node{
				Tag:     "website",
				Content: []byte(web),
			})
		}
	}

	bizNode := waBinary.Node{
		Tag: "business_profile",
		Attrs: waBinary.Attrs{
			"v":             "3",
			"mutation_type": "delta",
		},
		Content: nodes,
	}

	queryNode := waBinary.Node{
		Tag: "iq",
		Attrs: waBinary.Attrs{
			"to":    types.ServerJID.String(),
			"type":  "set",
			"xmlns": "w:biz",
		},
		Content: []waBinary.Node{bizNode},
	}

	return c.Client.DangerousInternals().SendNode(ctx, queryNode)
}
