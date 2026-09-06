package meowbail

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"os/exec"

	waBinary "go.mau.fi/whatsmeow/binary"
	"go.mau.fi/whatsmeow/types"
)

// SetGroupProfilePicture mengonversi, menormalkan, dan mengatur foto profil grup (avatar)
// Sesuai standar protokol resmi WhatsApp Web/Mobile:
// - Crop rasio 1:1 persegi (center crop)
// - Resolusi standar 640x640 piksel
// - Format Baseline JPEG dengan sub-sampling 4:2:0 murni tanpa metadata yang merusak
func (c *Client) SetGroupProfilePicture(ctx context.Context, groupJID types.JID, rawImageBytes []byte) (string, error) {
	if len(rawImageBytes) == 0 {
		// Hapus foto profil grup
		return c.SetGroupPhoto(ctx, groupJID, nil)
	}

	// 1. Normalisasi gambar menjadi baseline JPEG 640x640 menggunakan ffmpeg dengan kualitas sedang (quality ~50)
	// WhatsApp Web / Baileys standard: sharp(buffer).resize(640, 640).jpeg({ quality: 50 })
	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-i", "pipe:0",
		"-vf", "scale='if(gt(a,1),-1,640)':'if(gt(a,1),640,-1)',crop=640:640",
		"-pix_fmt", "yuvj420p",
		"-vcodec", "mjpeg",
		"-q:v", "7",
		"-f", "image2",
		"pipe:1",
	)
	cmd.Stdin = bytes.NewReader(rawImageBytes)
	jpegBytes, err := cmd.Output()

	// 2. Fallback jika ffmpeg gagal: gunakan image decoding Go standard library (quality 60)
	if err != nil || len(jpegBytes) == 0 {
		img, _, decErr := image.Decode(bytes.NewReader(rawImageBytes))
		if decErr != nil {
			return "", fmt.Errorf("format gambar tidak dikenali: %w", decErr)
		}

		b := img.Bounds()
		w := b.Dx()
		h := b.Dy()
		minDim := w
		if h < minDim {
			minDim = h
		}

		// Center crop
		startX := b.Min.X + (w-minDim)/2
		startY := b.Min.Y + (h-minDim)/2

		dst := image.NewRGBA(image.Rect(0, 0, 640, 640))
		for y := 0; y < 640; y++ {
			for x := 0; x < 640; x++ {
				srcX := startX + (x*minDim)/640
				srcY := startY + (y*minDim)/640
				dst.Set(x, y, img.At(srcX, srcY))
			}
		}

		var buf bytes.Buffer
		encErr := jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 60})
		if encErr != nil {
			return "", fmt.Errorf("gagal encode jpeg: %w", encErr)
		}
		jpegBytes = buf.Bytes()
	}

	// 3. Log detail ukuran JPEG sebelum dikirim
	// WhatsApp menolak profile picture jika payload melebih batas max IQ bytes
	return c.SetGroupPhoto(ctx, groupJID, jpegBytes)
}

// DownloadMediaDirectly mengunduh byte media dari stream pembaca WhatsApp
func (c *Client) DownloadMediaDirectly(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}

// EnsureGroupIQ Node query helper
func (c *Client) SendRawGroupNode(ctx context.Context, node waBinary.Node) error {
	return c.Client.DangerousInternals().SendNode(ctx, node)
}
