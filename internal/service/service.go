package service

import "taskapi/internal/repository"

type service struct {
	repo repository.Repositorier
}

// New ...
func New(repo repository.Repositorier) (Servicer, error) {
	return &service{
		repo: repo,
	}, nil
}
