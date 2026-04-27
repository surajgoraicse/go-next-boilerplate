package auth

import (
	"github.com/labstack/echo/v5"
	"github.com/surajgoraicse/go-next-boilerplate/internal/common/core"
)

func RegisterPublicRoutes(e *echo.Group, handler *Handler) {
	authRouter := e.Group("/v1/auth")
	authRouter.POST("/signin", core.WithBody(handler.SignIn))
	authRouter.POST("/register", core.WithBody(handler.Register))
	authRouter.POST("/refresh", core.WithBody(handler.Refresh))
}

func RegisterPrivateRoutes(e *echo.Group, handler *Handler, middleware echo.MiddlewareFunc) {
	authRouter := e.Group("/v1/auth")
	authRouter.Use(middleware)
	authRouter.POST("/logout", handler.Logout)
	authRouter.POST("/logout-all", handler.LogoutAll)
}
