package meowbail

import (
	"context"
	"io"
	"net/http"
	"os"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"
)

// DownloadMedia downloads media from a DownloadableMessage
func (c *Client) DownloadMedia(ctx context.Context, msg whatsmeow.DownloadableMessage) ([]byte, error) {
	if msg == nil {
		return nil, io.ErrNoProgress
	}
	return c.Client.Download(ctx, msg)
}

// UploadMedia uploads media and returns upload info
func (c *Client) UploadMedia(ctx context.Context, data []byte, mediaType whatsmeow.MediaType) (*whatsmeow.UploadResponse, error) {
	resp, err := c.Client.Upload(ctx, data, mediaType)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// FetchURL fetches data from a URL with timeout
func FetchURL(url string, timeout ...time.Duration) ([]byte, error) {
	t := 30 * time.Second
	if len(timeout) > 0 {
		t = timeout[0]
	}

	client := &http.Client{Timeout: t}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// ReadFile reads a file and returns its contents
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// SaveFile saves data to a file
func SaveFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

// BuildImageMessage builds an image message with newsletter context
func BuildImageMessage(resp *whatsmeow.UploadResponse, caption string, cfg *Config) *waE2E.Message {
	return &waE2E.Message{
		ImageMessage: &waE2E.ImageMessage{
			URL:           &resp.URL,
			Mimetype:      proto.String("image/jpeg"),
			Caption:       proto.String(caption),
			FileEncSHA256: resp.FileEncSHA256,
			FileSHA256:    resp.FileSHA256,
			FileLength:    &resp.FileLength,
			DirectPath:    &resp.DirectPath,
			MediaKey:      resp.MediaKey,
			ContextInfo:   buildNewsletterContext(cfg),
		},
	}
}

// BuildVideoMessage builds a video message with newsletter context
func BuildVideoMessage(resp *whatsmeow.UploadResponse, caption string, cfg *Config) *waE2E.Message {
	return &waE2E.Message{
		VideoMessage: &waE2E.VideoMessage{
			URL:           &resp.URL,
			Mimetype:      proto.String("video/mp4"),
			Caption:       proto.String(caption),
			FileEncSHA256: resp.FileEncSHA256,
			FileSHA256:    resp.FileSHA256,
			FileLength:    &resp.FileLength,
			DirectPath:    &resp.DirectPath,
			MediaKey:      resp.MediaKey,
			ContextInfo:   buildNewsletterContext(cfg),
		},
	}
}

// BuildDocumentMessage builds a document message with newsletter context
func BuildDocumentMessage(resp *whatsmeow.UploadResponse, filename, mimetype, caption string, cfg *Config) *waE2E.Message {
	return &waE2E.Message{
		DocumentMessage: &waE2E.DocumentMessage{
			URL:           &resp.URL,
			Mimetype:      proto.String(mimetype),
			FileName:      proto.String(filename),
			Caption:       proto.String(caption),
			FileEncSHA256: resp.FileEncSHA256,
			FileSHA256:    resp.FileSHA256,
			FileLength:    &resp.FileLength,
			DirectPath:    &resp.DirectPath,
			MediaKey:      resp.MediaKey,
			ContextInfo:   buildNewsletterContext(cfg),
		},
	}
}
