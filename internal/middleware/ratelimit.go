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
	visitors     sync.Map
	cleanupMu    sync.Mutex
	cleanupClose chan struct{}
)

// StartRateLimitCleanup launches the background goroutine that prunes stale
// visitor entries. Safe to call multiple times; only the first starts the loop.
func StartRateLimitCleanup() {
	cleanupMu.Lock()
	defer cleanupMu.Unlock()
	if cleanupClose != nil {
		return
	}
	cleanupClose = make(chan struct{})
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				visitors.Range(func(key, value any) bool {
					v := value.(*visitor)
					if time.Since(v.lastSeen) > 3*time.Minute {
						visitors.Delete(key)
					}
					return true
				})
			case <-cleanupClose:
				return
			}
		}
	}()
}

// StopRateLimitCleanup signals the cleanup goroutine to stop. Intended for
// use during graceful shutdown.
func StopRateLimitCleanup() {
	cleanupMu.Lock()
	defer cleanupMu.Unlock()
	if cleanupClose == nil {
		return
	}
	close(cleanupClose)
	cleanupClose = nil
}

func RateLimit(rps float64, burst int) fiber.Handler {
	return func(c fiber.Ctx) error {
		ip := c.IP()

		v, _ := visitors.LoadOrStore(ip, &visitor{
			limiter: rate.NewLimiter(rate.Limit(rps), burst),
		})
		vis := v.(*visitor)
		vis.lastSeen = time.Now()

		if !vis.limiter.Allow() {
			return fiber.NewError(fiber.StatusTooManyRequests, "rate limit exceeded")
		}
		return c.Next()
	}
}
