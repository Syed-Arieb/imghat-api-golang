package handler

import "github.com/gofiber/fiber/v3"

// Compress will compress a PNG or WebP image to a target quality level.
func Compress(c fiber.Ctx) error {
	return stub(c, "compress")
}
