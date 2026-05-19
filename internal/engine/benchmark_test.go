package engine

import (
	"image"
	"image/color"
	"image/png"
	"bytes"
	"testing"
)

func benchImage(w, h int) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetNRGBA(x, y, color.NRGBA{
				R: uint8(x * 4),
				G: uint8(y * 4),
				B: 128,
				A: 255,
			})
		}
	}
	return img
}

func benchPNG(w, h int) []byte {
	var buf bytes.Buffer
	png.Encode(&buf, benchImage(w, h))
	return buf.Bytes()
}

func BenchmarkCompress(b *testing.B) {
	data := benchPNG(800, 600)
	qualities := []float32{90, 60, 20}

	for _, q := range qualities {
		b.Run("", func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				Compress(data, CompressOptions{Quality: q})
			}
		})
	}
}

func BenchmarkResize(b *testing.B) {
	data := benchPNG(1920, 1080)
	modes := []ResizeMode{ResizeFit, ResizeFill, ResizeExact}

	for _, m := range modes {
		b.Run(string(m), func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				Resize(data, ResizeOptions{Width: 800, Height: 600, Mode: m, Quality: 80})
			}
		})
	}
}

func BenchmarkFit(b *testing.B) {
	data := benchPNG(800, 600)
	ratios := []struct { w, h int }{{16, 9}, {1, 1}, {4, 3}}

	for _, r := range ratios {
		b.Run("", func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				Fit(data, FitOptions{RatioW: r.w, RatioH: r.h, Quality: 80})
			}
		})
	}
}

func BenchmarkEncode(b *testing.B) {
	img := benchImage(800, 600)

	b.Run("PNG", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			Encode(img, EncodeOptions{Quality: 80, Format: FormatPNG})
		}
	})

	b.Run("WebP", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			Encode(img, EncodeOptions{Format: FormatWebP})
		}
	})
}
