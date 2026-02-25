package testdata

import (
	"database/sql"
	"fmt"
)

const (
	AppID     = 1
	AppSecret = "test_secret_key"
	AppName   = "test-app"

	AdminEmail    = "kefir.n@yandex.ru"
	AdminPassword = "Kefir4c_cc"
	AdminHash     = "$2a$10$Py02QUawy814cL6Nl3C7yuWe5fRT8N7ArPGoy0dxGvDi3OJiu2M6m"
)

func Seed(db *sql.DB) error {
	queries := []string{
		"TRUNCATE users, apps RESTART IDENTITY CASCADE;",

		fmt.Sprintf(`INSERT INTO apps (id, name, secret) VALUES
			(%d, '%s', '%s')`, AppID, AppName, AppSecret),

		fmt.Sprintf(`INSERT INTO users (email, pass_hash, is_admin) VALUES
			('%s','%s',true)`, AdminEmail, AdminHash),
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return fmt.Errorf("seed failed: %w", err)
		}
	}
	return nil
}
