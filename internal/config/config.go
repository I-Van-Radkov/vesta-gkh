package config

import (
	"fmt"
	"os"
	"time"

	newsUsecase "github.com/I-Van-Radkov/vesta-gkh/internal/usecase/news"
	postgres "github.com/I-Van-Radkov/vesta-gkh/pkg/db"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port         int           `env:"PORT" env-default:"8080"`
	ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"30s"`
	WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"30s"`

	GHTimeout time.Duration `env:"GRACEFUL_SHUTDOWN_TIMEOUT" env-default:"15s"`

	newsUsecase.ParserConfig

	postgres.PostgresConfig
}

func ParseConfigFromEnv() (*Config, error) {
	var cfg Config

	envPath := os.Getenv("ENV_PATH")
	if envPath == "" {
		envPath = "./config/.env"
	}

	if err := cleanenv.ReadConfig(envPath, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", envPath, err)
	}

	return &cfg, nil
}
