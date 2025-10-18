package service

import "github.com/tasks-control/core-back-end/internal/repository"

type Service struct {
	Repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{}
}
