package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string
	DBDSN   string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		AppPort: os.Getenv("APP_PORT"),
		DBDSN:   os.Getenv("DB_DSN"),
	}

	if cfg.AppPort == "" {
		cfg.AppPort = "9090"
	}

	if cfg.DBDSN == "" {
		return nil, fmt.Errorf("DB_DSN is not set")
	}

	return cfg, nil
}
