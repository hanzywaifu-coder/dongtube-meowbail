package meowbail

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hanzywaifu-coder/dongtube-meowbail/core"
)

// AddStickerMetadata menyuntikkan metadata EXIF resmi (PackName, Publisher, Emojis) ke dalam binary file WebP menggunakan webpmux
func AddStickerMetadata(webpBytes []byte, packName, publisher string) ([]byte, error) {
	if len(webpBytes) == 0 {
		return nil, fmt.Errorf("webp bytes kosong")
	}

	meta := core.StickerMetadata{
		PackName:  packName,
		Publisher: publisher,
	}
	exifChunk, err := core.BuildExifChunk(meta)
	if err != nil {
		return webpBytes, err
	}

	tmpWebp, err := os.CreateTemp("", "webp-*.webp")
	if err != nil {
		return webpBytes, nil
	}
	tmpExif, err := os.CreateTemp("", "exif-*.exif")
	if err != nil {
		_ = os.Remove(tmpWebp.Name())
		return webpBytes, nil
	}

	defer os.Remove(tmpWebp.Name())
	defer os.Remove(tmpExif.Name())

	_, _ = tmpWebp.Write(webpBytes)
	_, _ = tmpExif.Write(exifChunk)
	_ = tmpWebp.Close()
	_ = tmpExif.Close()

	outPath := tmpWebp.Name() + ".out.webp"
	defer os.Remove(outPath)

	cmd := exec.Command("webpmux", "-set", "exif", tmpExif.Name(), tmpWebp.Name(), "-o", outPath)
	if err := cmd.Run(); err != nil {
		// Fallback jika webpmux gagal
		return webpBytes, nil
	}

	muxedBytes, err := os.ReadFile(outPath)
	if err != nil || len(muxedBytes) == 0 {
		return webpBytes, nil
	}

	return muxedBytes, nil
}
