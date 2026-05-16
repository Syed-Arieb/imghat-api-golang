package engine

import (
	"fmt"
	"image"
	"image/color"

	"github.com/disintegration/imaging"
)

// Controls fit behaviour
type FitOptions struct {
	// RatioW and RatioH define target aspect ratio.
	RatioW  int
	RatioH  int
	Quality float32
	Format  Format // png or webp
}

// Pad img with transparent pixels to match target aspect ratio
func Fit(data []byte, opts FitOptions) ([]byte, error) {
	if opts.RatioW <= 0 || opts.RatioH <= 0 {
		return nil, fmt.Errorf("fit: invalid ratio %d:%d", opts.RatioW, opts.RatioH)
	}

	decoded, err := Decode(data)
	if err != nil {
		return nil, fmt.Errorf("fit: %w", err)
	}

	fitted := fitImage(decoded.Image, opts.RatioW, opts.RatioH)

	out, err := Encode(fitted, EncodeOptions{Quality: opts.Quality, Format: opts.Format})
	if err != nil {
		return nil, fmt.Errorf("fit: %w", err)
	}

	return out, nil
}

// Compute canvas size and paste the image centered.
func fitImage(img image.Image, ratioW, ratioH int) image.Image {
	bounds := img.Bounds()
	srcW, srcH := bounds.Dx(), bounds.Dy()

	// Calculate canvas dimensions that satisfy the ratio and contain the image.
	targetW, targetH := canvasSize(srcW, srcH, ratioW, ratioH)

	padX := (targetW - srcW) / 2
	padY := (targetH - srcH) / 2

	// Create a transparent NRGBA canvas.
	canvas := imaging.New(targetW, targetH, color.NRGBA{0, 0, 0, 0})

	// Paste the original image centered on the canvas.
	return imaging.Paste(canvas, img, image.Pt(padX, padY))
}

// Return the smallest canvas size that has exact aspect ratio and is enough to contain image
func canvasSize(srcW, srcH, ratioW, ratioH int) (int, int) {
	// Try fitting by width first
	targetH := srcW * ratioH / ratioW
	if targetH >= srcH {
		return srcW, targetH
	}
	// Height is the limiting dimension - fit by height
	targetW := srcH * ratioW / ratioH
	return targetW, srcH
}
