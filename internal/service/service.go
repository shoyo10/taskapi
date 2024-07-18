package service

import (
	"context"
	"taskapi/internal/model"
	"taskapi/internal/repository"
)

type service struct {
	repo repository.Repositorier
}

// New ...
func New(repo repository.Repositorier) Servicer {
	return &service{
		repo: repo,
	}
}

func (svc *service) CreateTask(ctx context.Context, task *model.Task) error {
	return svc.repo.CreateTask(ctx, task)
}

func (svc *service) ListTask(ctx context.Context) ([]model.Task, error) {
	return svc.repo.ListTask(ctx)
}

func (svc *service) UpdateTask(ctx context.Context, filter model.TaskFilter, in model.UpdateTaskInput) error {
	return svc.repo.UpdateTask(ctx, filter, in)
}

func (svc *service) DeleteTask(ctx context.Context, filter model.TaskFilter) error {
	return svc.repo.DeleteTask(ctx, filter)
}
