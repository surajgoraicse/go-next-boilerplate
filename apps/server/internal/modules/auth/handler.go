package auth

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v5"
	auth_mw "github.com/surajgoraicse/go-next-boilerplate/internal/common/middleware/auth"
	"github.com/surajgoraicse/go-next-boilerplate/internal/common/response"
	"github.com/surajgoraicse/go-next-boilerplate/internal/common/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Register(c *echo.Context, req RegisterRequest) error {
	statusCode, err := h.service.Register(c.Request().Context(), req)
	if err != nil {
		return response.NewResponse(c, statusCode, "registration failed", nil, err)
	}
	return response.NewResponse(c, statusCode, "registration successful, please verify email", nil, nil)
}

func (h *Handler) SignIn(c *echo.Context, req LoginRequest) error {
	// Extract device ID from header
	deviceIDStr := c.Request().Header.Get("X-Device-ID")
	var deviceID pgtype.UUID
	if deviceIDStr != "" {
		parsed, err := uuid.Parse(deviceIDStr)
		if err == nil {
			deviceID = pgtype.UUID{Bytes: parsed, Valid: true}
		}
	}

	tokens, sessionDeviceID, statusCode, err := h.service.SignIn(c.Request().Context(), req, deviceID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return response.NewResponse(c, statusCode, "user not found", nil, err)
		}
		if errors.Is(err, ErrEmailNotVerified) {
			return response.NewResponse(c, statusCode, "email not verified", nil, err)
		}
		if errors.Is(err, ErrInvalidCredentials) {
			return response.NewResponse(c, statusCode, "invalid credentials", nil, err)
		}
		return response.NewResponse(c, statusCode, "login failed", nil, err)
	}

	res := AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		DeviceID:     uuid.UUID(sessionDeviceID.Bytes).String(),
	}

	return response.NewResponse(c, statusCode, "login successful", res, nil)
}

func (h *Handler) Refresh(c *echo.Context, req RefreshRequest) error {
	tokens, statusCode, err := h.service.RefreshToken(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return response.NewResponse(c, statusCode, "refresh failed", nil, err)
	}

	return response.NewResponse(c, statusCode, "refresh successful", tokens, nil)
}

func (h *Handler) Logout(c *echo.Context) error {
	claims := auth_mw.GetUserClaims(c)
	deviceIDStr := c.Request().Header.Get("X-Device-ID")
	if deviceIDStr == "" {
		return response.NewResponse(c, http.StatusBadRequest, "device_id is required for logout", nil, nil)
	}

	deviceID, err := utils.StringToUUID(deviceIDStr)
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "invalid device_id", nil, err)
	}

	userID := pgtype.UUID{Bytes: claims.UserID, Valid: true}
	err = h.service.Logout(c.Request().Context(), userID, deviceID)
	if err != nil {
		return response.NewResponse(c, http.StatusInternalServerError, "logout failed", nil, err)
	}

	return response.NewResponse(c, http.StatusOK, "logged out successfully", nil, nil)
}

func (h *Handler) LogoutAll(c *echo.Context) error {
	claims := auth_mw.GetUserClaims(c)
	userID := pgtype.UUID{Bytes: claims.UserID, Valid: true}
	err := h.service.LogoutAll(c.Request().Context(), userID)
	if err != nil {
		return response.NewResponse(c, http.StatusInternalServerError, "logout all failed", nil, err)
	}

	return response.NewResponse(c, http.StatusOK, "logged out of all devices", nil, nil)
}
