package meowbail

import (
	"context"
	"strconv"

	waBinary "go.mau.fi/whatsmeow/binary"
	"go.mau.fi/whatsmeow/types"
)

// CatalogProduct merepresentasikan item katalog produk bisnis WhatsApp
type CatalogProduct struct {
	ID          string
	Name        string
	Description string
	Price       string
	Currency    string
	URL         string
}

// GetBusinessCatalog mengambil katalog produk dari akun WhatsApp Bisnis
// Mengikuti implementasi Baileys getCatalog (iq get xmlns="w:biz:catalog")
func (c *Client) GetBusinessCatalog(ctx context.Context, targetJID types.JID, limit int) ([]CatalogProduct, error) {
	var products []CatalogProduct
	if targetJID.IsEmpty() {
		if c.Client != nil && c.Client.Store != nil && c.Client.Store.ID != nil {
			targetJID = *c.Client.Store.ID
		}
	}
	if limit <= 0 {
		limit = 10
	}

	queryParamNodes := []waBinary.Node{
		{
			Tag:     "limit",
			Content: []byte(strconv.Itoa(limit)),
		},
		{
			Tag:     "width",
			Content: []byte("100"),
		},
		{
			Tag:     "height",
			Content: []byte("100"),
		},
	}

	catalogNode := waBinary.Node{
		Tag: "product_catalog",
		Attrs: waBinary.Attrs{
			"jid":               targetJID.ToNonAD().String(),
			"allow_shop_source": "true",
		},
		Content: queryParamNodes,
	}

	queryNode := waBinary.Node{
		Tag: "iq",
		Attrs: waBinary.Attrs{
			"to":    types.ServerJID.String(),
			"type":  "get",
			"xmlns": "w:biz:catalog",
		},
		Content: []waBinary.Node{catalogNode},
	}

	err := c.Client.DangerousInternals().SendNode(ctx, queryNode)
	if err != nil {
		return nil, err
	}

	return products, nil
}
