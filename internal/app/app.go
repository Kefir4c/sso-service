package app

import (
	"log/slog"

	grpcapp "github.com/Kefir4c/sso-service/internal/app/grpc"
	"github.com/Kefir4c/sso-service/internal/cache/redis"
	"github.com/Kefir4c/sso-service/internal/config"
	"github.com/Kefir4c/sso-service/internal/lib/logger/sl"
	auth2 "github.com/Kefir4c/sso-service/internal/services/auth"
	"github.com/Kefir4c/sso-service/internal/storage/postgres"
)

// App represents main application structure.
type App struct {
	GRPCServer *grpcapp.App
}

// New creates new application instance.
// Connects to PostgreSQL and Redis, creates auth service and gRPC server.
func New(log *slog.Logger, cfg *config.Config) *App {

	var (
		storage *postgres.Storage
		cache   *redis.Redis
		err     error
	)

	switch cfg.Storage.Type {
	case "postgres":
		storage, err = postgres.NewFromConfig(cfg)

		if err != nil {
			log.Error("failed to connect to Postgres", sl.Err(err))
			panic(err)
		}
	default:
		panic("unknow DB")
	}

	switch cfg.Cache.Driver {
	case "redis":
		cache, err = redis.NewFromConfig(cfg)

		if err != nil {
			log.Error("failed to connect to Redis", sl.Err(err))
			panic(err)
		}
	default:
		panic("unknow Cache")
	}

	auth := auth2.New(log, storage, storage, cache, cfg.GRPC.Timeout)

	gRPCServer := grpcapp.New(log, auth, cfg.GRPC.Port, cfg.GRPC.Timeout)

	return &App{
		GRPCServer: gRPCServer,
	}

}
