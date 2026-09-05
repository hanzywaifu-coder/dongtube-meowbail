package meowbail

import (
	"context"
	"log"

	"go.mau.fi/whatsmeow"
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

	return &Client{
		Client: cli,
		config: cfg,
	}
}

// NewClientFromWhatsmeow wraps an existing whatsmeow client
func NewClientFromWhatsmeow(cli *whatsmeow.Client, config ...*Config) *Client {
	cfg := DefaultConfig()
	if len(config) > 0 && config[0] != nil {
		cfg = config[0]
	}

	cli.EnableAutoReconnect = cfg.AutoReconnect

	return &Client{
		Client: cli,
		config: cfg,
	}
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

// Connect wraps whatsmeow's Connect
func (c *Client) Connect(ctx context.Context) error {
	return c.Client.Connect()
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
