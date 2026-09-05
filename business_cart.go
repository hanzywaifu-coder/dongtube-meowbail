package meowbail

import (
	"context"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// SendOrderReceipt mengirim struk belanja resmi pesanan (Order Receipt Message)
func (c *Client) SendOrderReceipt(ctx context.Context, chat types.JID, orderID, orderTitle string, itemCount int32, thumb []byte, sellerJID string) error {
	if sellerJID == "" {
		sellerJID = "0@s.whatsapp.net"
	}

	status := waE2E.OrderMessage_OrderStatus(1)
	surface := waE2E.OrderMessage_OrderSurface(1)

	msg := &waE2E.Message{
		OrderMessage: &waE2E.OrderMessage{
			OrderID:    proto.String(orderID),
			OrderTitle: proto.String(orderTitle),
			ItemCount:  proto.Int32(itemCount),
			Status:     &status,
			Surface:    &surface,
			SellerJID:  proto.String(sellerJID),
			Thumbnail:  thumb,
		},
	}

	_, err := c.Client.SendMessage(ctx, chat, msg)
	return err
}

// SendPaymentInvite mengirim ajakan transfer pembayaran (Payment Invite Message)
func (c *Client) SendPaymentInvite(ctx context.Context, chat types.JID, serviceType waE2E.PaymentInviteMessage_ServiceType, expiry int64) error {
	msg := &waE2E.Message{
		PaymentInviteMessage: &waE2E.PaymentInviteMessage{
			ServiceType:     serviceType.Enum(),
			ExpiryTimestamp: proto.Int64(expiry),
		},
	}
	_, err := c.Client.SendMessage(ctx, chat, msg)
	return err
}

// SendForwardedNewsletterPost meneruskan postingan dari channel WhatsApp ke chat/grup
func (c *Client) SendForwardedNewsletterPost(ctx context.Context, chat types.JID, text string, channelJID types.JID, channelName string, serverMessageID int64) error {
	msg := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(text),
			ContextInfo: &waE2E.ContextInfo{
				IsForwarded:     proto.Bool(true),
				ForwardingScore: proto.Uint32(9999),
				ForwardedNewsletterMessageInfo: &waE2E.ContextInfo_ForwardedNewsletterMessageInfo{
					NewsletterJID:   proto.String(channelJID.String()),
					NewsletterName:  proto.String(channelName),
					ServerMessageID: proto.Int32(int32(serverMessageID)),
				},
			},
		},
	}
	_, err := c.Client.SendMessage(ctx, chat, msg)
	return err
}
