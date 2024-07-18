package http

import (
	"github.com/labstack/echo/v4"
)

// SetRoutes ...
func SetRoutes(e *echo.Echo, h Handler) {
	taskAPI := e.Group("")
	taskAPI.GET("/tasks", h.ListTask)
	taskAPI.POST("/tasks", h.CreateTask)
	taskAPI.PUT("/tasks/:id", h.UpdateTask)
	taskAPI.DELETE("/tasks/:id", h.DeleteTask)
}
