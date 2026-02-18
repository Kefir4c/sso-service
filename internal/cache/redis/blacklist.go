package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const blacklistPrefix = "blacklist:"

func (r *Redis) AddToBlacklist(ctx context.Context, token string, ttl time.Duration) error {
	const op = "redise.AddToBlacklist "

	key := blacklistPrefix + token

	return r.client.Set(ctx, key, true, ttl).Err()
}

func (r *Redis) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	const op = "refis.IsBlacklisted"

	key := blacklistPrefix + token

	if err := r.client.Get(ctx, key).Err(); err != nil {
		if err == redis.Nil {

			return false, nil
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}
	return true, nil
}
