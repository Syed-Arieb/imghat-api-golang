package middleware

import (
	"io"

	"github.com/gofiber/fiber/v3"
)

const defaultMaxBytes = 10 << 20 // 10 MB

const ctxFileBytes = "file_bytes"

// Byte signatures used to detect format from raw bytes
var (
	pngMagic  = []byte{0x89, 0x50, 0x4E, 0x47}
	riffMagic = []byte("RIFF")
	webpSig   = []byte("WEBP")
)

// Middleware that rejects invalid uploads and stores the validated file bytes
// in the request context so the handler can reuse them without re-reading.
func ValidateImage(maxBytes int64) fiber.Handler {
	if maxBytes == 0 {
		maxBytes = defaultMaxBytes
	}

	return func(c fiber.Ctx) error {
		fh, err := c.FormFile("file")
		if err != nil {
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

		buf := make([]byte, fh.Size)
		if _, err := io.ReadFull(f, buf); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "could not read file")
		}

		if !isPNG(buf) && !isWebP(buf) {
			return fiber.NewError(fiber.StatusUnsupportedMediaType,
				"only PNG and WebP images are supported")
		}

		c.Locals(ctxFileBytes, buf)
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
