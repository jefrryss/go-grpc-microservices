package pg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

type Migrator struct {
	db  *sql.DB
	dir string
}

func New(db *sql.DB, dir string) *Migrator {
	return &Migrator{db: db, dir: dir}
}

func (m *Migrator) Up(ctx context.Context) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}
	if err := goose.UpContext(ctx, m.db, m.dir); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}
