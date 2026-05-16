package handler

import (
	"fmt"
	"imghat/internal/engine"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
)

// Fit handles POST /v1/image/fit
func Fit(c fiber.Ctx) error {
	data, err := readFile(c)
	if err != nil {
		return err
	}

	ratioW, ratioH, err := parseRatio(c.FormValue("ratio", ""))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	format := parseFormat(c.FormValue("format", "png"))

	out, err := engine.Fit(data, engine.FitOptions{
		RatioW:  ratioW,
		RatioH:  ratioH,
		Quality: parseQuality(c.FormValue("quality", "80")),
		Format:  engine.Format(format),
	})

	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	return imageResponse(c, out, format)
}

// "16:9" → (16, 9)
func parseRatio(s string) (int, int, error) {
	if s == "" {
		return 0, 0, fmt.Errorf("ratio is required (e.g. '1:1' or '16:9')")
	}

	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid ratio format %q, expected W:H", s)
	}

	w, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	h, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err1 != nil || err2 != nil || w <= 0 || h <= 0 {
		return 0, 0, fmt.Errorf("ratio values must be positive integers, got %q", s)
	}

	return w, h, nil
}
