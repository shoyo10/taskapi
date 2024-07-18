package http

import (
	"github.com/labstack/echo/v4"
)

// @title Task API Document
// @version 1.0
// @description This is task api document.

// @contact.name Shoyo
// @contact.url https://github.com/shoyo10/taskapi

// @host localhost:9090
// @BasePath /

// SetRoutes ...
func SetRoutes(e *echo.Echo, h Handler) {
	taskAPI := e.Group("")
	taskAPI.GET("/tasks", h.ListTask)
	taskAPI.POST("/tasks", h.CreateTask)
	taskAPI.PUT("/tasks/:id", h.UpdateTask)
	taskAPI.DELETE("/tasks/:id", h.DeleteTask)
}
