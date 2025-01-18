package config

import (
	"fmt"
	"github.com/PechatnovVladimir/concurrency_go/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type App struct {
	Name    string `envconfig:"APP_NAME"    required:"true" default:"kv-db"`
	Version string `envconfig:"APP_VERSION" required:"true" default:"v0.1.0"`
	Host    string `envconfig:"HOST" required:"true" default:"localhost:12345"`
}

type Config struct {
	App    App
	Logger logger.Config
}

func New() (Config, error) {
	var config Config

	err := godotenv.Load(".env")
	if err != nil {
		return config, fmt.Errorf("godotenv.Load: %w", err)
	}

	err = envconfig.Process("", &config)
	if err != nil {
		return config, fmt.Errorf("envconfig.Process: %w", err)
	}

	return config, nil
}
