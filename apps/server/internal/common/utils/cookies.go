package utils

import (
	"net/http"

	"github.com/labstack/echo/v5"
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
