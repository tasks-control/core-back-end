package main

import (
	"github.com/sirupsen/logrus"
	"github.com/tasks-control/core-back-end/internal/config"
	"github.com/tasks-control/core-back-end/internal/handler"
	"github.com/tasks-control/core-back-end/internal/repository"
	"github.com/tasks-control/core-back-end/internal/server"
	"github.com/tasks-control/core-back-end/internal/service"
	"github.com/tasks-control/core-back-end/pkg/utils"
)

func main() {
	log := utils.Logger()

	log.Info("Starting core-back-end")

	cfg := config.GetConfig()

	log.WithFields(logrus.Fields{
		"config": cfg,
	}).Debug("Config")

	repo, err := repository.New(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}

	svc := service.New(repo)
	h := handler.NewHandler(svc)

	server.NewServer(h).Run(cfg.ServerPort)
}
