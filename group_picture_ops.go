package meowbail

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"os/exec"
	"strings"

	waBinary "go.mau.fi/whatsmeow/binary"
	"go.mau.fi/whatsmeow/types"
)

// SetGroupProfilePicture mengonversi, menormalkan, dan mengatur foto profil grup (avatar)
// Menguji dual-method:
// Method 1: sendGroupIQ dengan Target = groupJID (standar whatsmeow & Baileys)
// Method 2: Fallback query direct to @g.us stanza jika server WhatsApp memerlukan format IQ set langsung
func (c *Client) SetGroupProfilePicture(ctx context.Context, groupJID types.JID, rawImageBytes []byte) (string, error) {
	if len(rawImageBytes) == 0 {
		return c.SetGroupPhoto(ctx, groupJID, nil)
	}

	// 1. Normalisasi gambar menjadi pure baseline JPEG 640x640 menggunakan ffmpeg
	// WhatsApp Web / Baileys: sharp(buffer).resize(640, 640).jpeg({ quality: 50 })
	// Menggunakan -fflags +bitexact -flags:v +bitexact untuk membersihkan komentar Lavc / metadata non-standar
	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-i", "pipe:0",
		"-vf", "scale=640:640:force_original_aspect_ratio=increase,crop=640:640",
		"-fflags", "+bitexact",
		"-flags:v", "+bitexact",
		"-pix_fmt", "yuvj420p",
		"-vcodec", "mjpeg",
		"-q:v", "6",
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

	// Method 1: Coba standard whatsmeow SetGroupPhoto
	picID, err := c.SetGroupPhoto(ctx, groupJID, jpegBytes)
	if err == nil && picID != "" {
		return picID, nil
	}

	// Method 2: Coba kirim via sendGroupIQ jika namespace w:profile:picture langsung
	groupNode := waBinary.Node{
		Tag:     "picture",
		Attrs:   waBinary.Attrs{"type": "image"},
		Content: jpegBytes,
	}
	resp, sendErr := c.Client.DangerousInternals().SendGroupIQ(ctx, "set", groupJID, groupNode)
	if sendErr == nil && resp != nil {
		picNode, ok := resp.GetOptionalChildByTag("picture")
		if ok {
			if id, okID := picNode.Attrs["id"].(string); okID && id != "" {
				return id, nil
			}
		}
		return "updated_group_iq", nil
	}

	if err != nil && (strings.Contains(err.Error(), "not a valid image") || strings.Contains(err.Error(), "not-acceptable")) {
		return "", fmt.Errorf("server WhatsApp menolak gambar (406 not-acceptable). Pastikan bot adalah admin dan grup mengizinkan anggota mengedit info grup")
	}

	return "", err
}
