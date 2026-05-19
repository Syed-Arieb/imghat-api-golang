package middleware

import (
	"bytes"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "imghat_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "imghat_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	requestInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "imghat_requests_in_flight",
			Help: "Current number of HTTP requests being handled",
		},
	)

	metricsOnce sync.Once
)

func initMetrics() {
	metricsOnce.Do(func() {
		prometheus.MustRegister(requestsTotal)
		prometheus.MustRegister(requestDuration)
		prometheus.MustRegister(requestInFlight)
		prometheus.MustRegister(collectors.NewGoCollector())
		prometheus.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	})
}

// MetricsMiddleware records request count, duration, and concurrent requests.
func MetricsMiddleware(c fiber.Ctx) error {
	initMetrics()
	start := time.Now()
	requestInFlight.Inc()
	defer requestInFlight.Dec()

	err := c.Next()

	status := strconv.Itoa(c.Response().StatusCode())
	method := c.Method()
	path := c.Route().Path

	requestsTotal.WithLabelValues(method, path, status).Inc()
	requestDuration.WithLabelValues(method, path).Observe(time.Since(start).Seconds())

	return err
}

// MetricsHandler serves the Prometheus metrics as plain text.
func MetricsHandler(c fiber.Ctx) error {
	initMetrics()
	var buf bytes.Buffer
	promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}).ServeHTTP(
		&fiberResponseWriter{buf: &buf}, newFiberRequest(c),
	)
	c.Set("Content-Type", "text/plain; version=0.0.4")
	return c.Send(buf.Bytes())
}

// fiberResponseWriter adapts a bytes.Buffer to satisfy http.ResponseWriter.
type fiberResponseWriter struct {
	buf *bytes.Buffer
}

func (w *fiberResponseWriter) Header() http.Header {
	return http.Header{}
}

func (w *fiberResponseWriter) Write(b []byte) (int, error) {
	return w.buf.Write(b)
}

func (w *fiberResponseWriter) WriteHeader(int) {}

// newFiberRequest creates a minimal http.Request for the promhttp handler.
func newFiberRequest(_ fiber.Ctx) *http.Request {
	return &http.Request{Method: "GET", RequestURI: "/metrics"}
}
