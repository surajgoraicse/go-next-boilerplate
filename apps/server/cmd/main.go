package main

import (
	"context"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	customMiddleware "github.com/surajgoraicse/go-next-boilerplate/internal/common/middleware"
	"github.com/surajgoraicse/go-next-boilerplate/internal/container"
	"github.com/surajgoraicse/go-next-boilerplate/internal/routes"
	"go.uber.org/zap"
)

// @title Swags API
// @version 0.1
// @description API documentation for swags.me
// @host localhost:8888
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize container
	di := container.NewContainer(ctx)
	defer di.Close()

	e := echo.New()
	allowedOrigins := []string{di.Config.FrontendOrigin}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           86400, // 24 hours
	}))
	e.Use(middleware.Recover())
	e.Use(customMiddleware.ZapLogger(di.Logger))
	e.Use(customMiddleware.PrometheusMiddleware())

	routes.RegisterRoutes(e, di)
	di.Logger.Info("Starting HTTP server", zap.String("port", di.Config.Port))

	if err := e.Start(":" + di.Config.Port); err != nil {
		di.Logger.Fatal("Failed to start server", zap.Error(err))
	}
}
