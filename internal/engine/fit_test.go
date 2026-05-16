package engine

import (
	"image"
	"image/color"
	"testing"
)

func TestCanvasSize(t *testing.T) {
	tests := []struct {
		name           string
		srcW, srcH     int
		ratioW, ratioH int
		wantW, wantH   int
	}{
		{"landscape→1:1", 30, 24, 1, 1, 30, 30},
		{"portrait→1:1", 24, 30, 1, 1, 30, 30},
		{"square→1:1", 30, 30, 1, 1, 30, 30},
		{"odd pixels→1:1", 31, 24, 1, 1, 31, 31},
		{"landscape→16:9", 160, 60, 16, 9, 160, 90},
		{"portrait→9:16", 90, 160, 9, 16, 90, 160},
		{"wide→4:3", 400, 100, 4, 3, 400, 300},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotW, gotH := canvasSize(tt.srcW, tt.srcH, tt.ratioW, tt.ratioH)
			if gotW != tt.wantW || gotH != tt.wantH {
				t.Errorf("canvasSize(%d,%d,%d:%d) = %d×%d, want %d×%d",
					tt.srcW, tt.srcH, tt.ratioW, tt.ratioH,
					gotW, gotH, tt.wantW, tt.wantH)
			}
		})
	}
}

func TestFitImageCentering(t *testing.T) {
	// 30×24 red image → 1:1 should produce a 30×30 canvas
	src := newSolidImage(30, 24, color.NRGBA{255, 0, 0, 255})
	result := fitImage(src, 1, 1)

	b := result.Bounds()
	if b.Dx() != 30 || b.Dy() != 30 {
		t.Fatalf("expected 30×30, got %d×%d", b.Dx(), b.Dy())
	}

	// Top-left pixel (0,0) must be transparent (padding row)
	r, g, bl, a := result.At(0, 0).RGBA()
	if a != 0 {
		t.Errorf("top-left should be transparent, got rgba(%d,%d,%d,%d)", r, g, bl, a)
	}

	// Center pixel must be red (original image)
	_, _, _, ca := result.At(15, 15).RGBA()
	if ca == 0 {
		t.Error("center pixel should not be transparent")
	}
}

func TestFitImageNoUpscale(t *testing.T) {
	// A square image fitted to 1:1 should not change at all
	src := newSolidImage(50, 50, color.NRGBA{0, 255, 0, 255})
	result := fitImage(src, 1, 1)
	b := result.Bounds()
	if b.Dx() != 50 || b.Dy() != 50 {
		t.Fatalf("expected 50×50, got %d×%d", b.Dx(), b.Dy())
	}
}

// Creates a colored image for tests.
func newSolidImage(w, h int, c color.NRGBA) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetNRGBA(x, y, c)
		}
	}
	return img
}
