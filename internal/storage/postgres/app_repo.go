package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Kefir4c/sso-service/internal/domain/models"
	"github.com/Kefir4c/sso-service/internal/storage"
	"github.com/jackc/pgx/v5"
)

const sqlGetApp = `
	SELECT id, name, secret
	FROM apps
	WHERE id = $1
`

func (s *Storage) App(ctx context.Context, appID int) (*models.App, error) {
	const op = "storage.app_repo.App"

	var app models.App
	if err := s.pool.QueryRow(ctx, sqlGetApp, appID).Scan(&app.ID, &app.Name, &app.Secret); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &app, nil
}
