package redis

import (
	"context"
	"fmt"
	"os"

	"github.com/Kefir4c/sso-service/internal/config"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewFromConfig(cfg *config.Config) (*Redis, error) {
	return New(fmt.Sprintf("%s:%d", cfg.Cache.Host, cfg.Cache.Port), os.Getenv("REDIS_PASS"), cfg.Cache.DB)
}

func New(addr, password string, db int) (*Redis, error) {
	const op = "redis.New"

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("%s: redis connect %w", op, err)
	}
	return &Redis{client: client}, nil
}
