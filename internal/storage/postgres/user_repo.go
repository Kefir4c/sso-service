package postgres

import (
	"context"
	"fmt"

	"github.com/Kefir4c/sso-service/internal/domain/models"
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
		pgErr, ok := err.(*pgconn.PgError)

		if ok && pgErr.Code == pgUniqueViolation {
			return 0, fmt.Errorf("%s: %w", op, err)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func(s *Storage) User(ctx context.Context,email string) (models.User,error) {
	const op = "Storage.user_repo.User"

	var user *models.User

	if err:= s.pool.QueryRow(ctx,sqlGetUserByEmail,email).
	Scan(&user.ID,&user.Email,&user.PassHash)
}
