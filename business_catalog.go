package meowbail

import (
	"context"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// SendCatalogCard mengirim katalog produk interaktif modern (WhatsApp Business Catalog Card)
func (c *Client) SendCatalogCard(ctx context.Context, chat types.JID, businessOwnerJID types.JID, title, description string, thumb []byte, productCount uint32) error {
	msg := &waE2E.Message{
		InteractiveMessage: &waE2E.InteractiveMessage{
			Header: &waE2E.InteractiveMessage_Header{
				HasMediaAttachment: proto.Bool(true),
				Media: &waE2E.InteractiveMessage_Header_ProductMessage{
					ProductMessage: &waE2E.ProductMessage{
						Catalog: &waE2E.ProductMessage_CatalogSnapshot{
							CatalogImage: &waE2E.ImageMessage{
								Mimetype:      proto.String("image/jpeg"),
								JPEGThumbnail: thumb,
							},
							Title:       proto.String(title),
							Description: proto.String(description),
						},
						BusinessOwnerJID: proto.String(businessOwnerJID.String()),
					},
				},
			},
			Body: &waE2E.InteractiveMessage_Body{
				Text: proto.String(title),
			},
			InteractiveMessage: &waE2E.InteractiveMessage_NativeFlowMessage_{
				NativeFlowMessage: &waE2E.InteractiveMessage_NativeFlowMessage{
					Buttons: []*waE2E.InteractiveMessage_NativeFlowMessage_NativeFlowButton{
						{
							Name:             proto.String("view_catalog"),
							ButtonParamsJSON: proto.String(`{}`),
						},
					},
				},
			},
		},
	}

	bizNodes := buildBizAdditionalNodes()
	_, err := c.Client.SendMessage(ctx, chat, msg, whatsmeow.SendRequestExtra{
		AdditionalNodes: &bizNodes,
	})
	return err
}
