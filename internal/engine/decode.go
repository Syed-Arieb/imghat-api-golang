package engine

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"

	"golang.org/x/image/webp"
)

type Format string

const (
	FormatPNG  Format = "png"
	FormatWebP Format = "webp"
)

type Decoded struct {
	Image  image.Image
	Format Format
}

var webpMagic = []byte("RIFF")

func Decode(data []byte) (*Decoded, error) {
	if isWebP(data) {
		img, err := webp.Decode(bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf("webp decode: %w", err)
		}
		return &Decoded{Image: img, Format: FormatWebP}, nil
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("png decode: %w", err)
	}
	return &Decoded{Image: img, Format: FormatPNG}, nil
}

func isWebP(data []byte) bool {
	return len(data) >= 12 &&
		bytes.Equal(data[0:4], webpMagic) &&
		bytes.Equal(data[8:12], []byte("WEBP"))
}
