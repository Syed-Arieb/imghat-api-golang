package handler

import "github.com/gofiber/fiber/v3"

func stub(c fiber.Ctx, op string) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"code":      fiber.StatusNotImplemented,
		"operation": op,
		"message":   op + " is not implemented yet",
	})
}
