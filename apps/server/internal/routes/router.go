package routes

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	authMiddleware "github.com/surajgoraicse/go-next-boilerplate/internal/common/middleware/auth"
	"github.com/surajgoraicse/go-next-boilerplate/internal/container"
)

// RegisterRoutes registers all the public and protected routes
func RegisterRoutes(e *echo.Echo, di *container.Container) {
	// health check api :
	e.GET("/health", healthCheck)

	apiRouter := e.Group("/api")
	protectedRouter := apiRouter.Group("")
	protectedRouter.Use(authMiddleware.AuthMiddleware(di.Config.JWTSecret))

	// Auth module routes (public and protected)
	// auth.RegisterPubicRoutes(apiRouter)

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
