package engine

import (
	"fmt"
	"image"

	"github.com/disintegration/imaging"
)

// Controls how the image fits target dimensions.
type ResizeMode string

const (
	// Scales the image to fit within w×h.
	ResizeFit ResizeMode = "fit"

	// Scales and center-crops to exactly w×h.
	ResizeFill ResizeMode = "fill"

	// Stretches to exactly w×h, ignoring aspect ratio.
	ResizeExact ResizeMode = "exact"
)

// Controls resize behaviour.
type ResizeOptions struct {
	Width     int
	Height    int
	Mode      ResizeMode
	FormatOut Format
	Quality   float32
}

// Decode -> Scale -> Re-Encode
func Resize(data []byte, opts ResizeOptions) ([]byte, Format, error) {
	decoded, err := Decode(data)
	if err != nil {
		return nil, "", fmt.Errorf("resize: %w", err)
	}

	resized := resizeImage(decoded.Image, opts)

	outFmt := opts.FormatOut
	if outFmt == "" {
		outFmt = decoded.Format
	}

	out, err := Encode(resized, EncodeOptions{Format: outFmt, Quality: opts.Quality})
	if err != nil {
		return nil, "", fmt.Errorf("resize: %w", err)
	}

	return out, outFmt, nil
}

func resizeImage(img image.Image, opts ResizeOptions) image.Image {
	bounds := img.Bounds()
	srcW, srcH := bounds.Dx(), bounds.Dy()

	// Never scale up
	if opts.Width >= srcW && opts.Height >= srcH {
		return img
	}

	switch opts.Mode {
	case ResizeFill:
		return imaging.Fill(img, opts.Width, opts.Height, imaging.Center, imaging.Lanczos)

	case ResizeExact:
		return imaging.Resize(img, opts.Width, opts.Height, imaging.Lanczos)

	default: // ResizeFit
		return imaging.Fit(img, opts.Width, opts.Height, imaging.Lanczos)
	}
}
