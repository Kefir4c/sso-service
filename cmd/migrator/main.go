package main

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/Kefir4c/sso-service/internal/config"
	"github.com/Kefir4c/sso-service/internal/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var steps int

	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)

	password := os.Getenv("POSTGRES_PASS")
	if password == "" {
		log.Error("POSTGRES_PASS environment variable is not set")
		panic("POSTGRES_PASS is required")
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Storage.Username,
		password,
		cfg.Storage.Host,
		cfg.Storage.Port,
		cfg.Storage.DBName,
	)
	migrationPath := cfg.MigrationsPath

	flag.IntVar(&steps, "steps", 0, "number of migration steps(positive for up, negative for down)")
	flag.Parse()

	if dbURL == "" {
		log.Error("database URL is empty")
		panic("database URL is required")
	}
	if migrationPath == "" {
		log.Error("migration path is empty")
		panic("migration path is required")
	}

	m, err := migrate.New("file://"+migrationPath, dbURL)
	if err != nil {
		log.Error("failed to create migration engine", slog.String("error", err.Error()))
		panic(err)
	}

	if steps != 0 {
		if err := m.Steps(steps); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				log.Info("no migrations to apply or rollback")
				return
			}
			log.Error("migration steps failed", slog.String("error", err.Error()))
			panic(err)
		}
		log.Info("migration steps completed", slog.Int("steps", steps))
		return
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("no new migrations to apply")
			return
		}
		log.Error("failed to apply migrations", slog.String("error", err.Error()))
		panic(err)
	}

	log.Info("migrations applied successfully")
}
