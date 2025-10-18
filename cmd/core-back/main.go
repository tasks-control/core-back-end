package main

import (
	"github.com/sirupsen/logrus"
	"github.com/tasks-control/core-back-end/internal/config"
	"github.com/tasks-control/core-back-end/internal/handler"
	"github.com/tasks-control/core-back-end/internal/server"
	"github.com/tasks-control/core-back-end/pkg/utils"
)

func main() {
	log := utils.Logger()

	log.Info("Starting core-back-end")

	cfg := config.GetConfig()

	log.WithFields(logrus.Fields{
		"config": cfg,
	}).Debug("Config")

	h := handler.NewHandler(nil)

	server.NewServer(h).Run(cfg.ServerPort)
}
