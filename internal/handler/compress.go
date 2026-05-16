package handler

import (
	"imghat/internal/engine"

	"github.com/gofiber/fiber/v3"
)

// Compress handles POST /v1/image/compress
func Compress(c fiber.Ctx) error {
	data, err := readFile(c)
	if err != nil {
		return err
	}

	quality := c.FormValue("quality", "80")

	out, err := engine.Compress(data, engine.CompressOptions{
		Quality: parseQuality(quality),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	return imageResponse(c, out)
}
