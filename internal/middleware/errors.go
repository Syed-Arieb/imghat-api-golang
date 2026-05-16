package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v3"
)

// Standard error structure returned by all endpoints.
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := "internal server error"

	if fe, ok := errors.AsType[*fiber.Error](err); ok {
		code = fe.Code
		msg = fe.Message
	} else if err != nil {
		msg = err.Error()
	}

	return c.Status(code).JSON(ErrorResponse{
		Code:    code,
		Message: msg,
	})
}
