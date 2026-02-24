package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	app2 "github.com/Kefir4c/sso-service/internal/app"
	"github.com/Kefir4c/sso-service/internal/config"
	"github.com/Kefir4c/sso-service/internal/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	log.Info("starting application", slog.Any("config", cfg))

	log.Info("success read config add setup logger")

	app := app2.New(log, cfg)

	go func() {
		app.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.GRPCServer.Stop()
	log.Info("Gracefull stopped")
}
