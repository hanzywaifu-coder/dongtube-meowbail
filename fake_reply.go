package meowbail

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func randHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// FakeQuotedBuilder penyedia aneka fake reply dari alipai-cmd.js / Evernight AI
type FakeQuotedBuilder struct{}

var FakeReply = &FakeQuotedBuilder{}

// Text membuat fake reply teks status broadcast (qtext / qtext2)
func (f *FakeQuotedBuilder) Text(text string) *waE2E.ContextInfo {
	return &waE2E.ContextInfo{
		StanzaID:      proto.String("DONGTUBE" + randHex(8)),
		Participant:   proto.String("0@s.whatsapp.net"),
		RemoteJID:     proto.String("status@broadcast"),
		QuotedMessage: &waE2E.Message{
			ExtendedTextMessage: &waE2E.ExtendedTextMessage{
				Text: proto.String(text),
			},
		},
	}
}

// Location membuat fake reply lokasi broadcast (qlocJpm / qlocPush)
func (f *FakeQuotedBuilder) Location(name string, thumb []byte) *waE2E.ContextInfo {
	return &waE2E.ContextInfo{
		StanzaID:      proto.String("DONGTUBE" + randHex(8)),
		Participant:   proto.String("0@s.whatsapp.net"),
		RemoteJID:     proto.String("status@broadcast"),
		QuotedMessage: &waE2E.Message{
			LocationMessage: &waE2E.LocationMessage{
				DegreesLatitude:  proto.Float64(0),
				DegreesLongitude: proto.Float64(0),
				Name:             proto.String(name),
				JPEGThumbnail:    thumb,
			},
		},
	}
}

// LiveLocation membuat fake reply live location broadcast (qlive)
func (f *FakeQuotedBuilder) LiveLocation(caption string, thumb []byte) *waE2E.ContextInfo {
	return &waE2E.ContextInfo{
		StanzaID:      proto.String("DONGTUBE" + randHex(8)),
		Participant:   proto.String("0@s.whatsapp.net"),
		RemoteJID:     proto.String("status@broadcast"),
		QuotedMessage: &waE2E.Message{
			LiveLocationMessage: &waE2E.LiveLocationMessage{
				Caption:       proto.String(caption),
				JPEGThumbnail: thumb,
			},
		},
	}
}

// Payment membuat fake reply payment request (qpayment)
func (f *FakeQuotedBuilder) Payment(botName string, amount1000 uint64, currency string) *waE2E.ContextInfo {
	if currency == "" {
		currency = "USD"
	}
	if amount1000 == 0 {
		amount1000 = 999999999
	}
	return &waE2E.ContextInfo{
		StanzaID:      proto.String("ownername"),
		Participant:   proto.String("0@s.whatsapp.net"),
		RemoteJID:     proto.String("0@s.whatsapp.net"),
		QuotedMessage: &waE2E.Message{
			RequestPaymentMessage: &waE2E.RequestPaymentMessage{
				CurrencyCodeIso4217: proto.String(currency),
				Amount1000:          proto.Uint64(amount1000),
				RequestFrom:         proto.String("0@s.whatsapp.net"),
				NoteMessage: &waE2E.Message{
					ExtendedTextMessage: &waE2E.ExtendedTextMessage{
						Text: proto.String(botName),
					},
				},
				ExpiryTimestamp: proto.Int64(999999999),
			},
		},
	}
}

// Toko membuat fake reply produk katalog marketplace (qtoko)
func (f *FakeQuotedBuilder) Toko(title, retailerId string, thumb []byte) *waE2E.ContextInfo {
	return &waE2E.ContextInfo{
		StanzaID:      proto.String("DONGTUBE" + randHex(8)),
		Participant:   proto.String("0@s.whatsapp.net"),
		RemoteJID:     proto.String("status@broadcast"),
		QuotedMessage: &waE2E.Message{
			ProductMessage: &waE2E.ProductMessage{
				Product: &waE2E.ProductMessage_ProductSnapshot{
					ProductImage: &waE2E.ImageMessage{
						Mimetype:      proto.String("image/jpeg"),
						JPEGThumbnail: thumb,
					},
					Title:             proto.String(title),
					CurrencyCode:      proto.String("IDR"),
					PriceAmount1000:   proto.Int64(999999999999999),
					RetailerID:        proto.String(retailerId),
					ProductImageCount: proto.Uint32(1),
				},
				BusinessOwnerJID: proto.String("0@s.whatsapp.net"),
			},
		},
	}
}

// Troli membuat fake reply troli belanja pengguna (troli)
func (f *FakeQuotedBuilder) Troli(senderJID types.JID, title, runtimeText string, thumb []byte) *waE2E.ContextInfo {
	participant := senderJID.ToNonAD().String()
	if participant == "" {
		participant = "0@s.whatsapp.net"
	}
	status := waE2E.OrderMessage_OrderStatus(1)
	surface := waE2E.OrderMessage_OrderSurface(1)
	return &waE2E.ContextInfo{
		StanzaID:      proto.String("DONGTUBE" + randHex(8)),
		Participant:   proto.String(participant),
		RemoteJID:     proto.String("status@broadcast"),
		QuotedMessage: &waE2E.Message{
			OrderMessage: &waE2E.OrderMessage{
				ItemCount:  proto.Int32(404),
				Status:     &status,
				Surface:    &surface,
				Message:    proto.String(runtimeText),
				OrderTitle: proto.String(title),
				SellerJID:  proto.String("0@s.whatsapp.net"),
				Thumbnail:  thumb,
				Token:      proto.String("AR6xBKmme9otv9WMZ4O6L9p968T2v99A"),
			},
		},
	}
}

// Kontak membuat fake reply kartu nama kontak developer (qkontak)
func (f *FakeQuotedBuilder) Kontak(displayName, phone string) *waE2E.ContextInfo {
	vcard := "BEGIN:VCARD\nVERSION:3.0\nFN:" + displayName + "\nORG:Developer\nTEL;type=CELL;type=VOICE;waid=" + phone + ":+" + phone + "\nEND:VCARD"
	return &waE2E.ContextInfo{
		StanzaID:      proto.String("DONGTUBE" + randHex(8)),
		Participant:   proto.String("0@s.whatsapp.net"),
		RemoteJID:     proto.String("status@broadcast"),
		QuotedMessage: &waE2E.Message{
			ContactMessage: &waE2E.ContactMessage{
				DisplayName: proto.String(displayName),
				Vcard:       proto.String(vcard),
			},
		},
	}
}

// CustomQuote membuat custom fake reply untuk target user apa pun
func (f *FakeQuotedBuilder) CustomQuote(remoteJID, participantJID string, msg *waE2E.Message) *waE2E.ContextInfo {
	return &waE2E.ContextInfo{
		StanzaID:      proto.String("DONGTUBE" + randHex(8)),
		Participant:   proto.String(participantJID),
		RemoteJID:     proto.String(remoteJID),
		QuotedMessage: msg,
	}
}

// SendTextWithFakeReply mengirim pesan teks menggunakan fake reply context
func (c *Client) SendTextWithFakeReply(ctx context.Context, chat types.JID, text string, fakeContext *waE2E.ContextInfo) error {
	if fakeContext == nil {
		fakeContext = &waE2E.ContextInfo{}
	}

	// Gabungkan newsletter context jika ada
	if c.config != nil && c.config.NewsletterJID != "" {
		fakeContext.IsForwarded = proto.Bool(true)
		fakeContext.ForwardingScore = proto.Uint32(9999)
		fakeContext.BusinessMessageForwardInfo = &waE2E.ContextInfo_BusinessMessageForwardInfo{
			BusinessOwnerJID: proto.String(c.config.BusinessOwnerJID),
		}
		fakeContext.ForwardedNewsletterMessageInfo = &waE2E.ContextInfo_ForwardedNewsletterMessageInfo{
			NewsletterJID:  proto.String(c.config.NewsletterJID),
			NewsletterName: proto.String(c.config.NewsletterName),
		}
	}

	msg := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text:        proto.String(text),
			ContextInfo: fakeContext,
		},
	}

	_, err := c.Client.SendMessage(ctx, chat, msg)
	return err
}
