package auth

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/surajgoraicse/go-next-boilerplate/internal/common/response"
	"github.com/surajgoraicse/go-next-boilerplate/internal/common/utils"
)

const (
	ClaimsKey = "claims"
)

// middleware to check if the user is authenticated
func AuthMiddleware(jwtSecret string) echo.MiddlewareFunc {
	echojwtConfig := echojwt.Config{
		SigningKey: []byte(jwtSecret),
		NewClaimsFunc: func(_ *echo.Context) jwt.Claims {
			return &utils.TokenPayload{}
		},
		// Prioritize header over cookie by listing header first
		TokenLookup: "header:Authorization:Bearer ,cookie:access_token",
		ErrorHandler: func(c *echo.Context, _ /* err */ error) error {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"message": "INVALID_TOKEN",
			})
		},
	}

	// create the wrapped middleware
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		// First apply Echo JWT middleware
		jwtMiddleware := echojwt.WithConfig(echojwtConfig)

		return jwtMiddleware(func(c *echo.Context) error {
			// Extract token from context (set by Echo JWT middleware)
			token, ok := c.Get("user").(*jwt.Token)
			if !ok {
				return response.NewResponse(c, http.StatusUnauthorized, "STATUS_UNAUTHORIZED", "INVALID_TOKEN_FORMAT", nil, nil)
			}

			// Parse claims into our custom struct
			claims, ok := token.Claims.(*utils.TokenPayload)
			if !ok {
				return response.NewResponse(c, http.StatusUnauthorized, "STATUS_UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
			}

			// Store parsed claims in context for easy access
			c.Set(ClaimsKey, claims)

			// Continue to next handler
			return next(c)
		})
	}
}

func GetUserClaims(c *echo.Context) *utils.TokenPayload {
	claims, ok := c.Get(ClaimsKey).(*utils.TokenPayload)
	if !ok {
		panic("User claims not found in context")
	}
	return claims
}
