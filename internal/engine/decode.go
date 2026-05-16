package engine

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"

	_ "github.com/HugoSmits86/nativewebp"
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

func Decode(data []byte) (*Decoded, error) {
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	if format != string(FormatPNG) && format != string(FormatWebP) {
		return nil, fmt.Errorf("decode: unsupported format %q", format)
	}
	return &Decoded{Image: img, Format: Format(format)}, nil
}
