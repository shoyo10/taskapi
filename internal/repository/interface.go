package repository

import (
	"context"
	"taskapi/internal/model"
)

type Repositorier interface {
	CreateTask(ctx context.Context, task *model.Task) error
	ListTask(ctx context.Context) ([]model.Task, error)
	UpdateTask(ctx context.Context, filter TaskFilter, in UpdateTaskInput) error
	DeleteTask(ctx context.Context, filter TaskFilter) error
}

type TaskFilter struct {
	ID int
}

type UpdateTaskInput struct {
	Name   *string
	Status *model.EnumTaskStatus
}
