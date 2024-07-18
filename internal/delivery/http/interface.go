package http

import (
	"github.com/labstack/echo/v4"
)

type Handler interface {
	// ListTask list task
	ListTask(c echo.Context) error
	// CreateTask create a task
	CreateTask(c echo.Context) error
	// UpdateTask update a task
	UpdateTask(c echo.Context) error
	// DeleteTask delete a task
	DeleteTask(c echo.Context) error
}
