package handler_test

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"

	"imghat/api"
	"imghat/internal/config"
	"imghat/internal/middleware"
)

func newApp() *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: middleware.ErrorHandler})
	api.RegisterRoutes(app, config.Load())
	return app
}

// makePNG generates a solid-color PNG of given dimensions in memory.
func makePNG(w, h int) []byte {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetNRGBA(x, y, color.NRGBA{100, 149, 237, 255})
		}
	}
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}

// multipartBody builds a multipart form with a file field plus extra fields.
func multipartBody(t *testing.T, fileData []byte, fields map[string]string) (*bytes.Buffer, string) {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, _ := w.CreateFormFile("file", "test.png")
	fw.Write(fileData)
	for k, v := range fields {
		w.WriteField(k, v)
	}
	w.Close()
	return &buf, w.FormDataContentType()
}

func TestCompressEndpoint(t *testing.T) {
	app := newApp()
	body, ct := multipartBody(t, makePNG(200, 200), map[string]string{"quality": "60"})

	req := httptest.NewRequest(http.MethodPost, "/v1/image/compress", body)
	req.Header.Set("Content-Type", ct)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "image/png" {
		t.Fatalf("expected image/png, got %s", ct)
	}
}

func TestResizeEndpoint(t *testing.T) {
	app := newApp()
	body, ct := multipartBody(t, makePNG(800, 600), map[string]string{
		"width": "400", "mode": "fit",
	})

	req := httptest.NewRequest(http.MethodPost, "/v1/image/resize", body)
	req.Header.Set("Content-Type", ct)

	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	// Decode response and check width was scaled down
	out, _ := io.ReadAll(resp.Body)
	img, _, _ := image.Decode(bytes.NewReader(out))
	if img.Bounds().Dx() > 400 {
		t.Errorf("expected width ≤ 400, got %d", img.Bounds().Dx())
	}
}

func TestFitEndpoint(t *testing.T) {
	app := newApp()
	// 30×24 → 1:1 should produce 30×30
	body, ct := multipartBody(t, makePNG(30, 24), map[string]string{"ratio": "1:1"})

	req := httptest.NewRequest(http.MethodPost, "/v1/image/fit", body)
	req.Header.Set("Content-Type", ct)

	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	out, _ := io.ReadAll(resp.Body)
	img, _, _ := image.Decode(bytes.NewReader(out))
	b := img.Bounds()
	if b.Dx() != 30 || b.Dy() != 30 {
		t.Errorf("expected 30×30, got %d×%d", b.Dx(), b.Dy())
	}
}

func TestValidationRejectsNonImage(t *testing.T) {
	app := newApp()
	body, ct := multipartBody(t, []byte("not an image"), nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/image/compress", body)
	req.Header.Set("Content-Type", ct)

	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusUnsupportedMediaType {
		t.Fatalf("expected 415, got %d", resp.StatusCode)
	}
}

func TestInvalidRatio(t *testing.T) {
	app := newApp()
	body, ct := multipartBody(t, makePNG(100, 100), map[string]string{"ratio": "bad"})

	req := httptest.NewRequest(http.MethodPost, "/v1/image/fit", body)
	req.Header.Set("Content-Type", ct)

	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}
