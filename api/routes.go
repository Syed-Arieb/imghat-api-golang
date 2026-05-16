package api

import (
	"github.com/gofiber/fiber/v3"

	"imghat/internal/config"
	"imghat/internal/handler"
	"imghat/internal/middleware"
)

// RegisterRoutes wires all versioned route groups onto the Fiber app.
func RegisterRoutes(app *fiber.App, cfg *config.Config) {
	v1 := app.Group("/v1")
	registerImageRoutes(v1, cfg)
}

func registerImageRoutes(r fiber.Router, cfg *config.Config) {
	img := r.Group("/image",
		middleware.RateLimit(cfg.RateLimit, cfg.RateBurst),
		middleware.ValidateImage(cfg.MaxFileSize),
	)

	// Accepts: multipart/form-data { file, quality (1-100), format (png|webp) }
	img.Post("/compress", handler.Compress)

	// Accepts: multipart/form-data { file, width, height, mode (fit|fill|exact) }
	img.Post("/resize", handler.Resize)

	// Accepts: multipart/form-data { file, ratio (e.g. "1:1"), format (png|webp) }
	img.Post("/fit", handler.Fit)
}
