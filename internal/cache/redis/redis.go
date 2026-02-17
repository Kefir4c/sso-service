package redis

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Kefir4c/sso-service/internal/config"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func New(addr, password string, db int) (*Redis, error) {
	const op = "redis.New"

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelCtx()

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("%s: redis connect %w", op, err)
	}
	return &Redis{client: client}, nil
}

func NewFromConfig(cfg *config.Config) (*Redis, error) {
	return New(fmt.Sprintf("%s:%d", cfg.Cache.Host, cfg.Cache.Port), os.Getenv("REDIS_PASS"), cfg.Cache.DB)
}

func (r *Redis) Close() error {
	const op = "redis.Close"

	if r.client == nil {
		return nil
	}

	if err := r.client.Close(); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	return nil
}
