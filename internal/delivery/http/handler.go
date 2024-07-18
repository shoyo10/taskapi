package http

import (
	"taskapi/internal/service"

	"github.com/labstack/echo/v4"
)

type handler struct {
	svc service.Servicer
}

// NewHandler create Handler instance
func NewHandler(svc service.Servicer) Handler {
	return &handler{
		svc: svc,
	}
}

func (h *handler) ListTask(c echo.Context) error {
	return nil
}

func (h *handler) CreateTask(c echo.Context) error {
	return nil
}

func (h *handler) UpdateTask(c echo.Context) error {
	return nil
}

func (h *handler) DeleteTask(c echo.Context) error {
	return nil
}
