package main

import (
	"log"

	"go.uber.org/zap"

	"github.com/sunr3d/basic-marketplace/internal/config"
	"github.com/sunr3d/basic-marketplace/internal/entrypoint"
	"github.com/sunr3d/basic-marketplace/internal/logger"
)

func main() {
	cfg, err := config.GetConfigFromEnv()
	if err != nil {
		log.Fatalf("config.GetConfigFromEnv: %s\n", err.Error())
	}

	zapLogger := logger.New(cfg.LogLevel)

	if err = entrypoint.Run(cfg, zapLogger); err != nil {
		zapLogger.Fatal("entrypoint.Run: ", zap.Error(err))
	}
}