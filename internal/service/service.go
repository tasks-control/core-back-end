package service

import (
	"errors"

	"github.com/tasks-control/core-back-end/internal/config"
	"github.com/tasks-control/core-back-end/internal/repository"
	"github.com/tasks-control/core-back-end/pkg/utils"
)

var (
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

type Service struct {
	Repo       repository.Repository
	JWTManager *utils.JWTManager
	JWTConfig  config.JWTConfig
}

func New(repo repository.Repository, cfg *config.Config) (*Service, error) {
	jwtManager, err := utils.NewJWTManager(cfg.JWT.SecretEnv)
	if err != nil {
		return nil, err
	}

	return &Service{
		Repo:       repo,
		JWTManager: jwtManager,
		JWTConfig:  cfg.JWT,
	}, nil
}
