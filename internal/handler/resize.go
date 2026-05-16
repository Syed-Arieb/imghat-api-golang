package handler

import (
	"imghat/internal/engine"

	"strconv"

	"github.com/gofiber/fiber/v3"
)

// Resize will scale an image down to the requested dimensions.
func Resize(c fiber.Ctx) error {
	data, err := readFile(c)
	if err != nil {
		return err
	}

	w, _ := strconv.Atoi(c.FormValue("width", "0"))
	h, _ := strconv.Atoi(c.FormValue("height", "0"))

	if w == 0 && h == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "at least one of width or height is required")
	}

	out, err := engine.Resize(data, engine.ResizeOptions{
		Width:   w,
		Height:  h,
		Mode:    engine.ResizeMode(c.FormValue("mode", "fit")),
		Quality: parseQuality(c.FormValue("quality", "80")),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	return imageResponse(c, out)
}
