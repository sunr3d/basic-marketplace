package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func GetConfigFromEnv() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("не удалось загрузить .env файл: \"%s\", используем дефолтные значения", err.Error())
	}

	cfg := &Config{}

	if err := envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("не удалось загрузить конфигурацию: %w", err)
	}

	return cfg, nil
}
