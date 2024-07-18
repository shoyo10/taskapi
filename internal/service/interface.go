package service

import (
	"context"
	"taskapi/internal/model"
)

type Servicer interface {
	CreateTask(ctx context.Context, task *model.Task) error
	ListTask(ctx context.Context) ([]model.Task, error)
	UpdateTask(ctx context.Context, filter model.TaskFilter, in model.UpdateTaskInput) error
	DeleteTask(ctx context.Context, filter model.TaskFilter) error
}
