package auth

import (
	"context"
	"time"

	"github.com/Kefir4c/sso-service/internal/domain/models"
)

type UserStorage interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (int64, error)
	User(ctx context.Context, email string) (*models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppStorage interface {
	App(ctx context.Context, id int) (*models.App, error)
}

type Cache interface {
	GetUser(ctx context.Context, email string) (*models.User, error)
	SetUser(ctx context.Context, user *models.User, ttl time.Duration) error
	AddToBlacklist(ctx context.Context, token string, ttl time.Duration) error
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}
