package meowbail

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// AlbumItem mendefinisikan 1 media dalam album (gambar atau video)
type AlbumItem struct {
	Image   []byte
	Video   []byte
	Caption string
}

// SendAlbum mengirim kumpulan foto/video dalam bentuk multi-media Album resmi WhatsApp
func (c *Client) SendAlbum(ctx context.Context, chat types.JID, items []AlbumItem) error {
	if len(items) < 2 {
		return fmt.Errorf("album requires at least 2 media items")
	}

	for _, itm := range items {
		if len(itm.Image) > 0 {
			resp, err := c.Client.Upload(ctx, itm.Image, whatsmeow.MediaImage)
			if err != nil {
				continue
			}
			msg := &waE2E.Message{
				ImageMessage: &waE2E.ImageMessage{
					URL:           &resp.URL,
					Mimetype:      proto.String("image/jpeg"),
					Caption:       proto.String(itm.Caption),
					FileEncSHA256: resp.FileEncSHA256,
					FileSHA256:    resp.FileSHA256,
					FileLength:    &resp.FileLength,
					DirectPath:    &resp.DirectPath,
					MediaKey:      resp.MediaKey,
					ContextInfo:   buildNewsletterContext(c.config),
				},
			}
			_, _ = c.Client.SendMessage(ctx, chat, msg)
		} else if len(itm.Video) > 0 {
			resp, err := c.Client.Upload(ctx, itm.Video, whatsmeow.MediaVideo)
			if err != nil {
				continue
			}
			msg := &waE2E.Message{
				VideoMessage: &waE2E.VideoMessage{
					URL:           &resp.URL,
					Mimetype:      proto.String("video/mp4"),
					Caption:       proto.String(itm.Caption),
					FileEncSHA256: resp.FileEncSHA256,
					FileSHA256:    resp.FileSHA256,
					FileLength:    &resp.FileLength,
					DirectPath:    &resp.DirectPath,
					MediaKey:      resp.MediaKey,
					ContextInfo:   buildNewsletterContext(c.config),
				},
			}
			_, _ = c.Client.SendMessage(ctx, chat, msg)
		}
	}

	return nil
}
