package engine

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
)

// Controls output Quality.
type EncodeOptions struct {
	Quality float32 // 1–100; maps to PNG compression level
}

// Encode encodes img as PNG.
func Encode(img image.Image, opts EncodeOptions) ([]byte, error) {
	var buf bytes.Buffer
	enc := &png.Encoder{CompressionLevel: qualityToPNGLevel(opts.Quality)}
	if err := enc.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("png encode: %w", err)
	}
	return buf.Bytes(), nil
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
