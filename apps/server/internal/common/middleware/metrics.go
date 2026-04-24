package middleware

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/surajgoraicse/go-next-boilerplate/internal/common/metrics"
)

// PrometheusMiddleware records per-request HTTP metrics into Prometheus.
// It uses c.Path() (the route pattern e.g. "/api/auth/:id") to avoid
// high cardinality from raw URL parameters.
func PrometheusMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			metrics.HTTPRequestsInFlight.Inc()
			start := time.Now()

			err := next(c)

			duration := time.Since(start).Seconds()

			// (*c).Response() returns http.ResponseWriter; the concrete type at
			// runtime is *echo.Response which exposes Status (int).
			var status int
			if resp, ok := (*c).Response().(*echo.Response); ok {
				status = resp.Status
			}
			if status == 0 {
				status = 200
			}

			method := (*c).Request().Method
			path := (*c).Path()
			if path == "" {
				path = (*c).Request().URL.Path
			}
			statusStr := strconv.Itoa(status)

			metrics.HTTPRequestsTotal.WithLabelValues(method, path, statusStr).Inc()
			metrics.HTTPRequestDuration.WithLabelValues(method, path, statusStr).Observe(duration)
			metrics.HTTPRequestsInFlight.Dec()

			return err
		}
	}
}
