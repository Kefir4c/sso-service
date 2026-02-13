package main

import (
	"log/slog"

	"github.com/Kefir4c/sso-service/internal/config"
	"github.com/Kefir4c/sso-service/internal/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	log.Info("starting application", slog.Any("config", cfg))

	// TODO: инициализировать логгер

	// TODO: инициализировать приложение (app)

	// TODO: запустить gRPC-сервер приложения
}
