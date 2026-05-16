package api

import (
	"github.com/gofiber/fiber/v3"

	"imghat/internal/handler"
)

// RegisterRoutes wires all versioned route groups onto the Fiber app.
func RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/v1")
	registerImageRoutes(v1)
}

func registerImageRoutes(r fiber.Router) {
	img := r.Group("/image")

	// POST /v1/image/compress
	// Accepts: multipart/form-data { file, quality (1-100), format (png|webp) }
	img.Post("/compress", handler.Compress)

	// POST /v1/image/resize
	// Accepts: multipart/form-data { file, width, height, mode (fit|fill|exact) }
	img.Post("/resize", handler.Resize)

	// POST /v1/image/fit
	// Accepts: multipart/form-data { file, ratio (e.g. "1:1"), format (png|webp) }
	img.Post("/fit", handler.Fit)
}
