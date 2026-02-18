package auth

import (
	"errors"
	"log/slog"
	"time"

	"github.com/Kefir4c/sso-service/internal/cache/redis"
	"github.com/Kefir4c/sso-service/internal/storage/postgres"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Auth struct {
	log      *slog.Logger
	storage  postgres.Storage
	cache    redis.Redis
	tokrnTTL time.Duration
}
