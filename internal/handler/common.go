package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// Extract raw bytes from the multipart "file" field.
// If ValidateImage middleware already read the file, the cached bytes
// stored in the request context are returned directly.
func readFile(c fiber.Ctx) ([]byte, error) {
	if data, ok := c.Locals("file_bytes").([]byte); ok {
		return data, nil
	}

	fh, err := c.FormFile("file")
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "missing file field")
	}

	f, err := fh.Open()
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "could not open upload")
	}

	defer f.Close()

	buf := make([]byte, fh.Size)
	if _, err := f.Read(buf); err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "could not read upload")
	}

	return buf, nil
}

// Write PNG bytes with the correct Content-Type.
func imageResponse(c fiber.Ctx, data []byte) error {
	c.Set(fiber.HeaderContentType, "image/png")
	return c.Send(data)
}

// Parses quality string to float32
func parseQuality(s string) float32 {
	v, err := strconv.ParseFloat(s, 32)

	if err != nil || v < 1 || v > 100 {
		return 80
	}

	return float32(v)
}
