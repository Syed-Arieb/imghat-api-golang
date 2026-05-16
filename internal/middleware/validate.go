package middleware

import (
	"github.com/gofiber/fiber/v3"
)

const defaultMaxBytes = 10 << 20 // 10 MB

// Byte signatures used to detect format from raw bytes
var (
	pngMagic  = []byte{0x89, 0x50, 0x4E, 0x47}
	riffMagic = []byte("RIFF")
	webpSig   = []byte("WEBP")
)

// Middleware that rejects invalid uploads.
func ValidateImage(maxBytes int64) fiber.Handler {
	if maxBytes == 0 {
		maxBytes = defaultMaxBytes
	}

	return func(c fiber.Ctx) error {
		fh, err := c.FormFile("file")
		if err != nil {
			// No file, let handler produce missing-field error.
			return c.Next()
		}

		if fh.Size > maxBytes {
			return fiber.NewError(fiber.StatusRequestEntityTooLarge,
				"file exceeds maximum allowed size")
		}

		f, err := fh.Open()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "could not inspect upload")
		}
		defer f.Close()

		header := make([]byte, 12)
		if _, err := f.Read(header); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "could not read file header")
		}

		if !isPNG(header) && !isWebP(header) {
			return fiber.NewError(fiber.StatusUnsupportedMediaType,
				"only PNG and WebP images are supported")
		}

		return c.Next()
	}
}

func isPNG(b []byte) bool {
	return len(b) >= 4 && string(b[0:4]) == string(pngMagic)
}

func isWebP(b []byte) bool {
	return len(b) >= 12 &&
		string(b[0:4]) == string(riffMagic) &&
		string(b[8:12]) == string(webpSig)
}
