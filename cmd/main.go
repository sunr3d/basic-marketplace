package main

import (
	"log"

	"go.uber.org/zap"

	"github.com/sunr3d/basic-marketplace/internal/config"
	"github.com/sunr3d/basic-marketplace/internal/entrypoint"
	"github.com/sunr3d/basic-marketplace/internal/logger"
)

func main() {
	// Загружаем конфиг из .env, либо дефолтные значения
	cfg, err := config.GetConfigFromEnv()
	if err != nil {
		log.Fatalf("config.GetConfigFromEnv: %s\n", err.Error())
	}

	// Создание логгера
	zapLogger := logger.New(cfg.LogLevel)

	// Точка входа в приложение
	if err = entrypoint.Run(cfg, zapLogger); err != nil {
		zapLogger.Fatal("entrypoint.Run: ", zap.Error(err))
	}
}
