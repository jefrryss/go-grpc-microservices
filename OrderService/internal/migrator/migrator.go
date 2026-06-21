package migrator

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

type Migrator struct {
	db             *sql.DB
	pathMigrations string
}

func NewMigrator(db *sql.DB, path string) *Migrator {
	return &Migrator{
		db:             db,
		pathMigrations: path,
	}
}

func (m *Migrator) Up() error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}

	if err := goose.Up(m.db, m.pathMigrations); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
