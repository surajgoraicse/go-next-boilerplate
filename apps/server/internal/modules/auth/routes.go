package auth

import (
	"github.com/labstack/echo/v5"
	"github.com/surajgoraicse/go-next-boilerplate/internal/common/core"
)

func RegisterPublicRoutes(e *echo.Group, handler *Handler) {
	authRouter := e.Group("/v1/auth")
	// authRouter.POST("/signin", core.WithBody(handler.SignIn))
	authRouter.POST("/register", core.WithBody(handler.Register))
	// authRouter.POST("/logout", handler.Logout)
	// authRouter.POST("/refresh", handler.Refresh)

	// Email verification
	// authRouter.POST("/verify-email", core.WithQuery(handler.VerifyEmail))
	// authRouter.POST("/verify-email/resend", core.WithBody(handler.SendVerificationEmail))

	// Password reset
	// authRouter.POST("/password/forgot", core.WithBody(handler.ForgotPassword))
	// authRouter.POST("/password/verify-otp", core.WithBody(handler.VerifyOTP))
	// authRouter.POST("/password/reset", core.WithBody(handler.ResetPassword))

	// Google OAuth routes
	// authRouter.GET("/google/login", handler.GoogleLogin)
	// authRouter.GET("/google/callback", handler.GoogleCallback)

}
