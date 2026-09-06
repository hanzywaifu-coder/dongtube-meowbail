package meowbail

import (
	"context"
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// ContactItem merepresentasikan satu kartu nama kontak
type ContactItem struct {
	Name  string
	Phone string
	Org   string
}

// SendContactsArray mengirim kumpulan kontak (ContactsArrayMessage) ke dalam obrolan
func (c *Client) SendContactsArray(ctx context.Context, chat types.JID, contacts []ContactItem) error {
	if len(contacts) == 0 {
		return fmt.Errorf("daftar kontak kosong")
	}

	var contactCards []*waE2E.ContactMessage
	displayName := contacts[0].Name
	if len(contacts) > 1 {
		displayName = fmt.Sprintf("%s dan %d kontak lainnya", contacts[0].Name, len(contacts)-1)
	}

	for _, contact := range contacts {
		cleanPhone := strings.TrimPrefix(strings.TrimSpace(contact.Phone), "+")
		vcardLines := []string{
			"BEGIN:VCARD",
			"VERSION:3.0",
			fmt.Sprintf("FN:%s", contact.Name),
			fmt.Sprintf("TEL;type=CELL;type=VOICE;waid=%s:+%s", cleanPhone, cleanPhone),
		}
		if contact.Org != "" {
			vcardLines = append(vcardLines, fmt.Sprintf("ORG:%s", contact.Org))
		}
		vcardLines = append(vcardLines, "END:VCARD")

		contactCards = append(contactCards, &waE2E.ContactMessage{
			DisplayName: proto.String(contact.Name),
			Vcard:       proto.String(strings.Join(vcardLines, "\n")),
		})
	}

	msg := &waE2E.Message{
		ContactsArrayMessage: &waE2E.ContactsArrayMessage{
			DisplayName: proto.String(displayName),
			Contacts:    contactCards,
			ContextInfo: buildNewsletterContext(c.config),
		},
	}

	_, err := c.Client.SendMessage(ctx, chat, msg)
	return err
}
