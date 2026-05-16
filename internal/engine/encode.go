package engine

import (
	"fmt"
	"image"
	"image/png"

	nativewebp "github.com/HugoSmits86/nativewebp"
)

// Controls output encoding.
type EncodeOptions struct {
	Quality float32 // 1–100;
	Format  Format  // png or webp
}

// Encode encodes img as PNG or WebP depending on opts.Format.
func Encode(img image.Image, opts EncodeOptions) ([]byte, error) {
	buf := getBuffer()
	defer putBuffer(buf)

	switch opts.Format {
	case FormatWebP:
		if err := nativewebp.Encode(buf, img, nil); err != nil {
			return nil, fmt.Errorf("webp encode: %w", err)
		}
	default:
		enc := &png.Encoder{CompressionLevel: qualityToPNGLevel(opts.Quality)}
		if err := enc.Encode(buf, img); err != nil {
			return nil, fmt.Errorf("png encode: %w", err)
		}
	}

	out := make([]byte, buf.Len())
	copy(out, buf.Bytes())
	return out, nil
}

func qualityToPNGLevel(q float32) png.CompressionLevel {
	switch {
	case q >= 80:
		return png.BestSpeed
	case q >= 40:
		return png.DefaultCompression
	default:
		return png.BestCompression
	}
}
