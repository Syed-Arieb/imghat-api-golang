package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	visitors = make(map[string]*visitor)
	mu       sync.Mutex
)

func init() {
	// Clean up stale visitors every minute
	go func() {
		for range time.Tick(time.Minute) {
			mu.Lock()
			for ip, v := range visitors {
				if time.Since(v.lastSeen) > 3*time.Minute {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()
}

func RateLimit(rps float64, burst int) fiber.Handler {
	return func(c fiber.Ctx) error {
		ip := c.IP()

		mu.Lock()
		v, ok := visitors[ip]
		if !ok {
			v = &visitor{limiter: rate.NewLimiter(rate.Limit(rps), burst)}
			visitors[ip] = v
		}
		v.lastSeen = time.Now()
		lim := v.limiter
		mu.Unlock()

		if !lim.Allow() {
			return fiber.NewError(fiber.StatusTooManyRequests, "rate limit exceeded")
		}
		return c.Next()
	}
}
