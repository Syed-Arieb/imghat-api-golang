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
	Width   int
	Height  int
	Mode    ResizeMode
	Quality float32
	Format  Format // png or webp
}

// Decode -> Scale -> Re-Encode
func Resize(data []byte, opts ResizeOptions) ([]byte, error) {
	decoded, err := Decode(data)
	if err != nil {
		return nil, fmt.Errorf("resize: %w", err)
	}

	resized := resizeImage(decoded.Image, opts)

	out, err := Encode(resized, EncodeOptions{Quality: opts.Quality, Format: opts.Format})
	if err != nil {
		return nil, fmt.Errorf("resize: %w", err)
	}

	return out, nil
}

func selectFilter(quality float32) imaging.ResampleFilter {
	switch {
	case quality >= 80:
		return imaging.Lanczos
	case quality >= 40:
		return imaging.CatmullRom
	default:
		return imaging.Linear
	}
}

func resizeImage(img image.Image, opts ResizeOptions) image.Image {
	bounds := img.Bounds()
	srcW, srcH := bounds.Dx(), bounds.Dy()

	w, h := opts.Width, opts.Height

	// Never scale up
	if (w == 0 || w >= srcW) && (h == 0 || h >= srcH) {
		return img
	}

	filter := selectFilter(opts.Quality)

	switch opts.Mode {
	case ResizeFill:
		if w == 0 {
			w = srcW
		}
		if h == 0 {
			h = srcH
		}
		return imaging.Fill(img, w, h, imaging.Center, filter)

	case ResizeExact:
		if w == 0 {
			w = srcW
		}
		if h == 0 {
			h = srcH
		}
		return imaging.Resize(img, w, h, filter)

	default: // ResizeFit
		if w == 0 || h == 0 {
			return imaging.Resize(img, w, h, filter)
		}
		return imaging.Fit(img, w, h, filter)
	}
}
