package engine

import (
	"bytes"
	"fmt"
	"image"
	"image/png"

	"github.com/chai2010/webp"
)

// EncodeOptions controls output format and quality
type EncodeOptions struct {
	Format  Format
	Quality float32
}

// Encode encodes img into a byte slice
func Encode(img image.Image, opts EncodeOptions) ([]byte, error) {
	var buf bytes.Buffer

	switch opts.Format {
	case FormatWebP:
		q := clampQuality(opts.Quality)
		if err := webp.Encode(&buf, img, &webp.Options{Lossless: false, Quality: q}); err != nil {
			return nil, fmt.Errorf("webp encode: %w", err)
		}

	case FormatPNG, "":
		// PNG compression level maps from quality
		enc := &png.Encoder{CompressionLevel: qualityToPNGLevel(opts.Quality)}
		if err := enc.Encode(&buf, img); err != nil {
			return nil, fmt.Errorf("png encode: %w", err)
		}

	default:
		return nil, fmt.Errorf("unsupported format: %s", opts.Format)
	}

	return buf.Bytes(), nil
}

func clampQuality(q float32) float32 {
	if q <= 0 {
		return 80 // sensible default
	}
	if q > 100 {
		return 100
	}
	return q
}

func qualityToPNGLevel(q float32) png.CompressionLevel {
	if q >= 80 {
		return png.BestSpeed
	}
	if q >= 40 {
		return png.DefaultCompression
	}
	return png.BestCompression
}
