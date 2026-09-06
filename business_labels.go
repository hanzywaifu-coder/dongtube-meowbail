package meowbail

import (
	"context"

	"go.mau.fi/whatsmeow/appstate"
	"go.mau.fi/whatsmeow/types"
)

// LabelChat memberikan atau mencabut label bisnis pada obrolan (WhatsApp Business Labels)
func (c *Client) LabelChat(ctx context.Context, chat types.JID, labelID string, labeled bool) error {
	patch := appstate.BuildLabelChat(chat, labelID, labeled)
	return c.Client.SendAppState(ctx, patch)
}

// LabelMessage memberikan atau mencabut label bisnis pada pesan spesifik
func (c *Client) LabelMessage(ctx context.Context, chat types.JID, labelID, messageID string, labeled bool) error {
	patch := appstate.BuildLabelMessage(chat, labelID, messageID, labeled)
	return c.Client.SendAppState(ctx, patch)
}

// EditLabel membuat, mengubah, atau menghapus label bisnis (nama & warna ID 0-19)
func (c *Client) EditLabel(ctx context.Context, labelID, labelName string, colorID int32, deleted bool) error {
	patch := appstate.BuildLabelEdit(labelID, labelName, colorID, deleted)
	return c.Client.SendAppState(ctx, patch)
}
