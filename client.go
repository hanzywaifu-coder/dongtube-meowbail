package meowbail

import (
	"context"
	"log"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

// NewClient creates a new dongtube-meowbail client
func NewClient(device interface{}, logger interface{}, config ...*Config) *Client {
	cfg := DefaultConfig()
	if len(config) > 0 && config[0] != nil {
		cfg = config[0]
	}

	// Use whatsmeow client as base
	var cli *whatsmeow.Client
	if device != nil {
		// Type assert to whatsmeow device store
		if ds, ok := device.(*whatsmeow.Client); ok {
			cli = ds
		}
	}

	if cli == nil {
		log.Fatal("dongtube-meowbail: invalid device store")
	}

	cli.EnableAutoReconnect = cfg.AutoReconnect

	store := NewMemoryMessageStore(5000)
	c := &Client{
		Client:       cli,
		config:       cfg,
		LIDResolver:  NewLIDResolver(),
		RetryTracker: NewRetrySpiralingTracker(5),
		Store:        store,
	}
	c.AttachMessageStore(store)
	return c
}

// NewClientFromWhatsmeow wraps an existing whatsmeow client
func NewClientFromWhatsmeow(cli *whatsmeow.Client, config ...*Config) *Client {
	cfg := DefaultConfig()
	if len(config) > 0 && config[0] != nil {
		cfg = config[0]
	}

	cli.EnableAutoReconnect = cfg.AutoReconnect

	store := NewMemoryMessageStore(5000)
	c := &Client{
		Client:       cli,
		config:       cfg,
		LIDResolver:  NewLIDResolver(),
		RetryTracker: NewRetrySpiralingTracker(5),
		Store:        store,
	}
	c.AttachMessageStore(store)
	return c
}

// SetNewsletter sets the newsletter context for all messages
func (c *Client) SetNewsletter(jid, name string) {
	c.config.NewsletterJID = jid
	c.config.NewsletterName = name
}

// SetBusinessOwner sets the business owner JID
func (c *Client) SetBusinessOwner(jid string) {
	c.config.BusinessOwnerJID = jid
}

// AddEventHandler wraps whatsmeow's AddEventHandler
func (c *Client) AddEventHandler(handler func(evt interface{})) uint32 {
	return c.Client.AddEventHandler(handler)
}

// Connect wraps whatsmeow's Connect and performs auto-follow if configured
func (c *Client) Connect(ctx context.Context) error {
	// Sync versi WhatsApp Web terbaru secara otomatis untuk mencegah 405 out of date
	go func() {
		ver, err := whatsmeow.GetLatestVersion(context.Background(), nil)
		if err == nil && ver != nil {
			store.SetWAVersion(*ver)
		}
	}()

	err := c.Client.Connect()
	if err != nil {
		return err
	}

	// Auto-follow channels in background once connected and logged in
	go func() {
		for i := 0; i < 30; i++ {
			if c.Client.IsLoggedIn() {
				break
			}
			time.Sleep(1 * time.Second)
		}

		if !c.Client.IsLoggedIn() {
			return
		}

		var jids []string
		if c.config.NewsletterJID != "" {
			jids = append(jids, c.config.NewsletterJID)
		}
		for _, j := range c.config.AutoFollowJIDs {
			jids = append(jids, j)
		}

		seen := make(map[string]bool)
		for _, rawJID := range jids {
			if rawJID == "" || seen[rawJID] {
				continue
			}
			seen[rawJID] = true
			parsed, err := types.ParseJID(rawJID)
			if err != nil {
				continue
			}
			_ = c.Client.FollowNewsletter(context.Background(), parsed)
		}
	}()

	return nil
}

// Disconnect wraps whatsmeow's Disconnect
func (c *Client) Disconnect() {
	c.Client.Disconnect()
}

// IsConnected wraps whatsmeow's IsConnected
func (c *Client) IsConnected() bool {
	return c.Client.IsConnected()
}

// IsLoggedIn wraps whatsmeow's IsLoggedIn
func (c *Client) IsLoggedIn() bool {
	return c.Client.IsLoggedIn()
}

// ParseMessageEvent parses a raw event into a MessageEvent
func ParseMessageEvent(evt interface{}) *MessageEvent {
	e, ok := evt.(*events.Message)
	if !ok {
		return nil
	}

	text := ""
	msg := e.Message
	if msg.GetConversation() != "" {
		text = msg.GetConversation()
	} else if msg.ExtendedTextMessage != nil {
		text = msg.ExtendedTextMessage.GetText()
	} else if msg.ImageMessage != nil {
		text = msg.ImageMessage.GetCaption()
	} else if msg.VideoMessage != nil {
		text = msg.VideoMessage.GetCaption()
	} else if msg.DocumentMessage != nil {
		text = msg.DocumentMessage.GetCaption()
	} else if msg.AudioMessage != nil {
		text = "" // audio has no text
	} else if msg.StickerMessage != nil {
		text = "" // sticker has no text
	}

	return &MessageEvent{
		Message:  e,
		Sender:   e.Info.Sender,
		Chat:     e.Info.Chat,
		Text:     text,
		IsGroup:  e.Info.Chat.Server == "g.us",
		IsFromMe: e.Info.IsFromMe,
	}
}

// GetConfig returns the client configuration
func (c *Client) GetConfig() *Config {
	return c.config
}
