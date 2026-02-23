package app

import (
	"log/slog"

	grpcapp "github.com/Kefir4c/sso-service/internal/app/grpc"
	"github.com/Kefir4c/sso-service/internal/cache/redis"
	"github.com/Kefir4c/sso-service/internal/config"
	"github.com/Kefir4c/sso-service/internal/lib/logger/sl"
	"github.com/Kefir4c/sso-service/internal/storage/postgres"
)

type App struct {
	gRPCServer *grpcapp.App
}

func New(log *slog.Logger,cfg *config.Config) *App{
	
	var(
		storage postgres.Storage
		cahce redis.Redis
		err error
	)

	switch cfg.Storage.Type{
	case "postgres":
		store,err:= postgres.NewFromConfig(cfg)

		if err != nil{
			log.Error("failed to cennect to Postgres",sl.Err(err))
			panic(err)
		}
	default:
		panic("unrnow DB")
	}

	switch cfg.Cache.Driver{
	case "redis":
		cache,err:= redis.NewFromConfig(cfg)
		
		if err != nil{
			log.Error("failed to connect to Redis",sl.Err(err))
			panic(err)
		}
	default:
		panic("unknow Cache")	
	}

	auth:= auth.

}
