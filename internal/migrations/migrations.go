package migrations

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS

func RunMigrations(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)
	err := goose.Up(db, ".")
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
