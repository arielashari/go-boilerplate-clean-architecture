package middleware

import (
	"fmt"
	"time"

	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/pkg/metrics"
	"github.com/gofiber/fiber/v3"
)

func MetricsMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		metrics.HttpRequestsInFlight.Inc()
		defer metrics.HttpRequestsInFlight.Dec()

		err := c.Next()

		status := fmt.Sprintf("%d", c.Response().StatusCode())
		path := c.Route().Path
		method := c.Method()
		duration := time.Since(start).Seconds()

		metrics.HttpRequestsTotal.WithLabelValues(method, path, status).Inc()
		metrics.HttpRequestDuration.WithLabelValues(method, path, status).Observe(duration)

		return err
	}
}
