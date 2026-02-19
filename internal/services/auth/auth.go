package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Kefir4c/sso-service/internal/cache/redis"
	"github.com/Kefir4c/sso-service/internal/domain/models"
	"github.com/Kefir4c/sso-service/internal/lib/logger/sl"
	"github.com/Kefir4c/sso-service/internal/storage"
	"github.com/Kefir4c/sso-service/internal/storage/postgres"
	"github.com/Kefir4c/sso-service/internal/validation"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrInvalidPassword    = errors.New("invalid password format")
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrAppNotFound        = errors.New("app not found")
	ErrTokenInvalid       = errors.New("token is invalid")
	ErrTokenBlacklisted   = errors.New("token is blacklisted")
)

const (
	cacheTTL         = 24 * time.Hour
	bcryptCost       = bcrypt.DefaultCost
	operationTimeout = 5 * time.Second
	bgTimeout        = 2 * time.Second
)

type Auth struct {
	log      *slog.Logger
	storage  postgres.Storage
	cache    redis.Redis
	tokenTTL time.Duration
}

func New(log *slog.Logger, storage postgres.Storage, cache redis.Redis, tokenTTL time.Duration) (*Auth, error) {
	return &Auth{
		log:      log,
		storage:  storage,
		cache:    cache,
		tokenTTL: tokenTTL,
	}, nil
}

func (a *Auth) getUser(ctx context.Context, email string) (*models.User, error) {
	ctx, cancelCtx := context.WithTimeout(ctx, operationTimeout)
	defer cancelCtx()

	log := a.log.With(slog.String("op", "auth.getUser"), slog.String("email", email))

	user, err := a.cache.GetUser(ctx, email)
	if err == nil {
		log.Debug("cache hit")
		return user, nil
	}
	if !errors.Is(err, redis.ErrUserNotFound) {
		log.Error("cache error", sl.Err(err))
		return nil, fmt.Errorf("cache error: %w", err)
	}

	log.Debug("cache miss")

	user, err = a.storage.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, fmt.Errorf("%w", ErrUserNotFound)
		}
		log.Error("storage error", sl.Err(err))
		return nil, fmt.Errorf("get user from storage: %w", err)
	}

	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), bgTimeout)
		defer cancel()
		if err := a.cache.SetUser(bgCtx, user, cacheTTL); err != nil {
			log.Error("failed to save user to cache", sl.Err(err))
		}
	}()
	return user, nil
}

func (a *Auth) Register(ctx context.Context, email, password string) (int64, error) {
	const op = "auth.Register"

	ctx, cancelCtx := context.WithTimeout(ctx, operationTimeout)
	defer cancelCtx()

	log := a.log.With(slog.String("op", op), slog.String("email", email))

	if err := validation.ValidateEmail(email); err != nil {
		log.Warn("invalid email", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidEmail)
	}

	if err := validation.ValidatePassword(password); err != nil {
		log.Warn("invalid password", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidPassword)
	}

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.storage.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists")
			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}
		log.Error("failed to save user to storage", op, err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), bgTimeout)
		defer cancel()

		user := &models.User{
			ID:       id,
			Email:    email,
			PassHash: passHash,
		}

		if err := a.cache.SetUser(bgCtx, user, cacheTTL); err != nil {
			log.Error("failed to save user to cache", sl.Err(err))
		}
	}()

	log.Info("user registered", slog.Int64("id", id))
	return id, nil
}
