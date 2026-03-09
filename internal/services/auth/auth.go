package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Kefir4c/sso-service/internal/domain/models"
	"github.com/Kefir4c/sso-service/internal/lib/jwt"
	"github.com/Kefir4c/sso-service/internal/lib/logger/sl"
	"github.com/Kefir4c/sso-service/internal/storage"
	"github.com/Kefir4c/sso-service/internal/validation"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
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
	log        *slog.Logger
	storage    UserStorage
	appStorage AppStorage
	cache      Cache
	tokenTTL   time.Duration
}

func New(log *slog.Logger,
	userStorage UserStorage,
	appStorage AppStorage,
	cache Cache,
	tokenTTL time.Duration) *Auth {
	return &Auth{
		log:        log,
		storage:    userStorage,
		appStorage: appStorage,
		cache:      cache,
		tokenTTL:   tokenTTL,
	}
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
	if !errors.Is(err, ErrUserNotFound) {
		log.Error("cache error", sl.Err(err))
		return nil, fmt.Errorf("cache error: %w", err)
	}

	log.Debug("cache miss")

	user, err = a.storage.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			fmt.Fprintf(os.Stderr, "🔍 DEBUG getUser: storage.ErrUserNotFound caught, returning ErrInvalidCredentials\n")
			return nil, ErrInvalidCredentials
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
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if err := validation.ValidatePassword(password); err != nil {
		log.Warn("invalid password", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.storage.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrLoginExists) {
			log.Warn("user already exists")
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
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

func (a *Auth) Login(ctx context.Context, email, password string, appID int) (string, error) {
	const op = "auth.Login"

	ctx, cancelCtx := context.WithTimeout(ctx, operationTimeout)
	defer cancelCtx()

	log := a.log.With(slog.String("op", op), slog.String("email", email))

	if err := validation.ValidateEmail(email); err != nil {
		log.Warn("invalid email", sl.Err(err))
		return "", fmt.Errorf("%s :%w", op, err)
	}

	if err := validation.ValidatePassword(password); err != nil {
		log.Warn("invalid password", sl.Err(err))
		return "", fmt.Errorf("%s :%w", op, err)
	}

	if appID < 0 {
		log.Warn("invalid app_id", slog.Int("app_id", appID))
		return "", fmt.Errorf("%s: invalid app_id", op)
	}

	log.Info("logging in user")

	user, err := a.getUser(ctx, email)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			return "", ErrInvalidCredentials
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Warn("invalid password")
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appStorage.App(ctx, appID)
	if err != nil {
		log.Error("failed to get app", slog.Int("app_id", appID), sl.Err(err))
		if errors.Is(err, ErrAppNotFound) {
			log.Warn("app not found", slog.Int("app_id", appID))
			return "", fmt.Errorf("%s :%w", op, ErrAppNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("failed to generate jwt", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in seccesfully")
	return token, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.IsAdmin"

	ctx, cancelCtx := context.WithTimeout(ctx, operationTimeout)
	defer cancelCtx()

	log := a.log.With(slog.String("op", op), slog.Int64("user_id", userID))

	if userID <= 0 {
		log.Warn("invalid user_id")
		return false, fmt.Errorf("%s : invalid user_id", op)
	}

	log.Info("checking admin status")

	isAdmin, err := a.storage.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found")
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		log.Error("failed to check admin status", sl.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("admin status check", slog.Bool("is_admin", isAdmin))
	return isAdmin, nil
}

func (a *Auth) ValidateToken(ctx context.Context, token string) (bool, int64, string, int) {
	const op = "auth.ValidateToken"

	ctx, cancelCtx := context.WithTimeout(ctx, operationTimeout)
	defer cancelCtx()

	log := a.log.With(slog.String("op", op))

	if token == "" {
		log.Warn("empty token")
		return false, 0, "", 0
	}

	blackListed, err := a.cache.IsBlacklisted(ctx, token)
	if err != nil {
		log.Error("failed to check blacklist", sl.Err(err))
		return false, 0, "", 0
	}
	if blackListed {
		log.Info("token is blacklisted")
		return false, 0, "", 0
	}

	claims, err := jwt.ParseToken(token)
	if err != nil {
		log.Info("failed to parse token", sl.Err(err))
		return false, 0, "", 0
	}

	app, err := a.appStorage.App(ctx, claims.AppID)
	if err != nil {
		log.Info("app not found", slog.Int("app_id", claims.AppID))
		return false, 0, "", 0
	}

	validatedClaims, err := jwt.ValidateTokenWithSecret(token, app.Secret)
	if err != nil {
		log.Info("invalid token signature", sl.Err(err))
		return false, 0, "", 0
	}

	_, err = a.storage.User(ctx, validatedClaims.Email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Info("user not found", slog.String("email", validatedClaims.Email))
			return false, 0, "", 0
		}
		log.Error("failed to check user existence", sl.Err(err))
		return false, 0, "", 0
	}

	log.Info("token validated successfully",
		slog.Int64("user_id", validatedClaims.UserID),
		slog.String("email", validatedClaims.Email),
		slog.Int("app_id", validatedClaims.AppID),
	)

	return true, validatedClaims.UserID, validatedClaims.Email, validatedClaims.AppID
}

func (a *Auth) Logout(ctx context.Context, token string) (bool, error) {
	const op = "auth.Logout"

	ctx, cancelCtx := context.WithTimeout(ctx, operationTimeout)
	defer cancelCtx()

	log := a.log.With(slog.String("op", op))

	if token == "" {
		log.Warn("empty token")
		return false, fmt.Errorf("%s: token is empty", op)
	}
	log.Info("Processing logout ")

	claims, err := jwt.ParseToken(token)
	if err != nil {
		log.Warn("invalid token in logout request", sl.Err(err))
		return false, nil
	}

	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		log.Info("token already expired")
		return true, nil
	}

	if err := a.cache.AddToBlacklist(ctx, token, ttl); err != nil {
		log.Error("failed to add token to blacklist", sl.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfully", slog.Duration("token_ttl", ttl))

	return true, nil
}
