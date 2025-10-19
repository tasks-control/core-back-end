package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/tasks-control/core-back-end/internal/repository"

	"github.com/go-playground/validator/v10"

	"github.com/tasks-control/core-back-end/pkg/utils"
	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	ServerPort string            `validate:"required" yaml:"serverPort"`
	Database   repository.Config `validate:"required" yaml:"database"`
	JWT        JWTConfig         `validate:"required" yaml:"jwt"`
}

type JWTConfig struct {
	SecretEnv            string `validate:"required" yaml:"secretEnv"`
	AccessTokenDuration  int    `validate:"required,min=60" yaml:"accessTokenDuration"`    // seconds
	RefreshTokenDuration int    `validate:"required,min=3600" yaml:"refreshTokenDuration"` // seconds
}

func GetConfig() (cfg *Config) {
	log := utils.Logger()
	configPath := flag.String("c", "./cmd/core-back/config.yaml", "path to config")
	flag.Parse()

	cfg = &Config{}

	err := read(*configPath, cfg)
	if err != nil {
		log.WithError(err).Fatal("can't read config")
	}

	v := validator.New()

	err = v.Struct(cfg)
	if err != nil {
		log.WithError(err).Fatal("can't validate config")
	}

	return cfg
}

func read(path string, cfg interface{}) error {
	data, err := os.ReadFile(path) //nolint:gosec // Config file path is controlled by application
	if err != nil {
		return fmt.Errorf("cant read config file: %s", err.Error())
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return fmt.Errorf("cant parse config: %s", err.Error())
	}

	return nil
}
