package engine

import "fmt"

// CompressOptions controls how the image is compressed.
type CompressOptions struct {
	Quality float32 // 1–100
}

// Compress decodes raw image bytes, re-encodes at the target quality
func Compress(data []byte, opts CompressOptions) ([]byte, error) {
	decoded, err := Decode(data)
	if err != nil {
		return nil, fmt.Errorf("compress: %w", err)
	}
	out, err := Encode(decoded.Image, EncodeOptions{Quality: opts.Quality})
	if err != nil {
		return nil, fmt.Errorf("compress: %w", err)
	}
	return out, nil
}
