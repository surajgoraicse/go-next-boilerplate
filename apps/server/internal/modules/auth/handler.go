package auth

import (
	"github.com/labstack/echo/v5"
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
	// statusCode, err := h.service.Register(c.Request().Context(), req)
	// if err != nil {
	// 	return err
	// }
	// return response.NewResponse(c, statusCode, map[string]interface{}{
	// 	"message": "registration successful, please verify email",
	// })
	return nil
}
