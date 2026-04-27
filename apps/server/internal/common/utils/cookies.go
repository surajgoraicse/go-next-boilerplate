package utils

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/surajgoraicse/go-next-boilerplate/internal/config"
)

// CookieConfig holds cookie configuration
type CookieConfig struct {
	Domain   string
	Secure   bool
	HTTPOnly bool
	SameSite http.SameSite
	Path     string
}

// DefaultCookieConfig returns default cookie configuration
func DefaultCookieConfig() CookieConfig {
	return CookieConfig{
		Domain:   "",
		Secure:   true,
		HTTPOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
}

// SetAccessTokenCookie sets the access token cookie with secure flags
func SetAccessTokenCookie(c *echo.Context, token string, maxAge int, config CookieConfig) {
	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     config.Path,
		Domain:   config.Domain,
		MaxAge:   maxAge,
		Secure:   config.Secure,
		HttpOnly: config.HTTPOnly,
		SameSite: config.SameSite,
	}
	c.SetCookie(cookie)
}

// SetRefreshTokenCookie sets the refresh token cookie with secure flags
func SetRefreshTokenCookie(c *echo.Context, token string, maxAge int, config CookieConfig) {
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Path:     config.Path,
		Domain:   config.Domain,
		MaxAge:   maxAge,
		Secure:   config.Secure,
		HttpOnly: config.HTTPOnly,
		SameSite: config.SameSite,
	}
	c.SetCookie(cookie)
}

// ClearAuthCookies clears both access and refresh token cookies
func ClearAuthCookies(c *echo.Context, config CookieConfig) {
	// Clear access token
	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     config.Path,
		Domain:   config.Domain,
		MaxAge:   -1,
		Secure:   config.Secure,
		HttpOnly: config.HTTPOnly,
		SameSite: config.SameSite,
	}
	c.SetCookie(accessCookie)

	// Clear refresh token
	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     config.Path,
		Domain:   config.Domain,
		MaxAge:   -1,
		Secure:   config.Secure,
		HttpOnly: config.HTTPOnly,
		SameSite: config.SameSite,
	}
	c.SetCookie(refreshCookie)
}

func SetAuthCookies(c *echo.Context, tokens AuthTokens, config *config.Config) {
	accessMaxAge := calculateCookieMaxAge(config.AccessTokenExpiry, 86400)
	accessCookie := NewSecureCookie("auth_token", tokens.AccessToken, accessMaxAge, "/")
	c.SetCookie(accessCookie)

	refreshMaxAge := h.calculateCookieMaxAge(h.service.config.RefreshTokenExpiry, 604800)
	refreshCookie := h.NewSecureCookie("refresh_token", tokens.RefreshToken, refreshMaxAge, "/")
	c.SetCookie(refreshCookie)
}

func calculateCookieMaxAge(expiry string, fallback int) int {
	duration, err := time.ParseDuration(expiry)
	if err != nil {
		return fallback
	}
	return int(duration.Seconds())
}

func NewSecureCookie(name string, value string, maxAge int, path string, appEnv config.Environment) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		HttpOnly: true,
		Secure:   appEnv == config.Production,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
	}
}
