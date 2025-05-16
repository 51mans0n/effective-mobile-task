// Package migration handles database schema migrations.
package migration

import (
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
)

// Во время сборки возьмёт все .sql файлы из папки sql/ и впечатает (embed) их прямо в бинарник
//
//go:embed sql/*.sql
var sqlFiles embed.FS

// Up applies all embedded SQL migrations to the connected database.
func Up(db *sqlx.DB) error {
	d, err := iofs.New(sqlFiles, "sql")
	if err != nil {
		return fmt.Errorf("create iofs: %w", err)
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
