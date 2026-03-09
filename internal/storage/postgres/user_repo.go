package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Kefir4c/sso-service/internal/domain/models"
	"github.com/Kefir4c/sso-service/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	sqlInsertUser = `
		INSERT INTO users (email, pass_hash)
		VALUES ($1, $2)
		RETURNING id
	`
	sqlGetUserByEmail = `
		SELECT id, email, pass_hash
		FROM users
		WHERE email = $1
	`
	sqlCheckAdmin = `
		SELECT is_admin
		FROM users
		WHERE id = $1
	`
	pgUniqueViolation = "23505"
)

func (s *Storage) SaveUser(ctx context.Context, email string, passhash []byte) (int64, error) {
	const op = "storage.user_repo.SaveUser"

	var id int64
	if err := s.pool.QueryRow(ctx, sqlInsertUser, email, passhash).
		Scan(&id); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrLoginExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (*models.User, error) {
	const op = "storage.user_repo.User"

	var user models.User
	if err := s.pool.QueryRow(ctx, sqlGetUserByEmail, email).
		Scan(&user.ID, &user.Email, &user.PassHash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.user_repo.IsAdmin"

	var isAdmin bool

	if err := s.pool.QueryRow(ctx, sqlCheckAdmin, userID).Scan(&isAdmin); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return isAdmin, nil
}
