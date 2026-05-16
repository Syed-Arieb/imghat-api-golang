package handler

import "github.com/gofiber/fiber/v3"

// Resize will scale an image down to the requested dimensions.
func Resize(c fiber.Ctx) error {
	return stub(c, "resize")
}
