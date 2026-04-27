package routes

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	authMiddleware "github.com/surajgoraicse/go-next-boilerplate/internal/common/middleware/auth"
	"github.com/surajgoraicse/go-next-boilerplate/internal/container"
	"github.com/surajgoraicse/go-next-boilerplate/internal/modules/auth"
	_ "github.com/surajgoraicse/go-next-boilerplate/swagger"
	echoSwagger "github.com/swaggo/echo-swagger/v2"
)

// RegisterRoutes registers all the public and protected routes
func RegisterRoutes(e *echo.Echo, di *container.Container) {
	e.GET("/health", healthCheck)

	// Prometheus metrics endpoint — not behind auth middleware
	e.GET("/metrics", metricsHandler)

	// Swagger docs
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	apiRouter := e.Group("/api")
	protectedRouter := e.Group("/api")
	protectedRouter.Use(authMiddleware.AuthMiddleware(di.Config.JwtSecret))

	// register auth routes
	auth.RegisterPublicRoutes(apiRouter, di.AuthHandler)
	auth.RegisterPrivateRoutes(protectedRouter, di.AuthHandler, authMiddleware.AuthMiddleware(di.Config.JwtSecret))

}

// metricsHandler adapts the standard http.Handler from promhttp to Echo.
func metricsHandler(c *echo.Context) error {
	h := promhttp.Handler()
	h.ServeHTTP((*c).Response(), (*c).Request())
	return nil
}

// healthCheck godoc
// @Summary Health check
// @Description Check if the server is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func healthCheck(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
