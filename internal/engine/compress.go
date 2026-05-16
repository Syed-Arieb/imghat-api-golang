package engine

import "fmt"

// CompressOptions controls how the image is compressed.
type CompressOptions struct {
	Quality   float32 // 1–100
	FormatOut Format
}

// Compress decodes raw image bytes, re-encodes at the target quality
func Compress(data []byte, opts CompressOptions) ([]byte, Format, error) {
	decoded, err := Decode(data)
	if err != nil {
		return nil, "", fmt.Errorf("compress: %w", err)
	}

	outFmt := opts.FormatOut
	if outFmt == "" {
		outFmt = decoded.Format // preserve source format
	}

	out, err := Encode(decoded.Image, EncodeOptions{
		Format:  outFmt,
		Quality: opts.Quality,
	})
	if err != nil {
		return nil, "", fmt.Errorf("compress: %w", err)
	}

	return out, outFmt, nil
}
