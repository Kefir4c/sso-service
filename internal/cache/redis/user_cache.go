package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Kefir4c/sso-service/internal/domain/models"
	"github.com/redis/go-redis/v9"
)

const (
	userPrefix = "user:"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

func (r *Redis) SetUser(ctx context.Context, user *models.User, ttl time.Duration) error {
	const op = "redis.SetUser"

	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	key := userPrefix + user.Email

	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *Redis) GetUser(ctx context.Context, email string) (*models.User, error) {
	const op = "redis.Getuser"

	key := userPrefix + email

	data, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var user models.User

	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	return &user, nil
}
