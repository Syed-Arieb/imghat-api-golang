package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"

	"github.com/goccy/go-json"

	"imghat/api"
	"imghat/internal/config"
	"imghat/internal/middleware"
)

const version = "0.1.0"

func main() {
	app := fiber.New(fiber.Config{
		AppName:      "imghat " + version,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		ErrorHandler: middleware.ErrorHandler,
	})

	// Global middleware
	app.Use(recover.New()) // catch panics, keep the server alive
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} ${method} ${path} ${latency}\n",
	}))

	// Health check - used by Docker / k8s probes
	app.Get("/healthz", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"version": version,
		})
	})

	cfg := config.Load()

	// All versioned image routes
	api.RegisterRoutes(app, cfg)

	// Start in background so we can listen for signals
	app.Listen(":" + cfg.Port)

	// Block until SIGINT / SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("shutting down gracefully…")
	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
	log.Println("bye")
}
