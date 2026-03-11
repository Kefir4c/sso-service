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

	UserEmail    = "Smetan.k@yandex.ru"
	UserPassword = "Smetan#2026#"
	UserHash     = "$2a$10$YNDAwyX7rCSNgQCwuexdVOwEHNMI4J.HvEqZzRxmWmgOdjKAjMxFm"
)

// Seed populates database with test data.
// Truncates existing tables, inserts test app and users.
func Seed(db *sql.DB) error {
	queries := []string{
		"TRUNCATE users, apps RESTART IDENTITY CASCADE;",

		fmt.Sprintf(`INSERT INTO apps (id, name, secret) VALUES
			(%d, '%s', '%s')`, AppID, AppName, AppSecret),

		fmt.Sprintf(`INSERT INTO users (email, pass_hash, is_admin) VALUES
			('%s','%s',true)`, AdminEmail, AdminHash),

		fmt.Sprintf(`INSERT INTO users (email,pass_hash,is_admin) VALUES 
    		('%s','%s',false)`, UserEmail, UserHash),
	}

	for _, q := range queries {
		fmt.Printf("DEBUG: executing query:%s\n", q)
		if _, err := db.Exec(q); err != nil {
			return fmt.Errorf("seed failed: %w", err)
		}
	}
	fmt.Println("DEBUG: seed completed successfully")
	return nil
}
