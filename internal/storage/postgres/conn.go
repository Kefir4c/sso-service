package postgres

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Kefir4c/sso-service/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, host, database, username, password string, port int) (*Storage, error) {
	const op = "postgres.New"

	connString := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		host, port, database, username, password)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf(" %s: parse config %w", op, err)
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("%s: connect %w", op, err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("%s: ping %w", op, err)
	}
	return &Storage{pool: pool}, nil
}

func NewFromConfig(config *config.Config) (*Storage, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelCtx()
	return New(
		ctx,
		config.Storage.Host,
		config.Storage.DBName,
		config.Storage.Username,
		os.Getenv("POSTGRES_PASS"),
		config.Storage.Port,
	)
}
