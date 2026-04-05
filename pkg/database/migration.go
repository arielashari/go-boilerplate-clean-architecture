package database

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func RunMigrations(pool *pgxpool.Pool) error {
	db := stdlib.OpenDB(*pool.Config().ConnConfig)

	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		return fmt.Errorf("migration driver error: %w", err)
	}

	subFS, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create sub-filesystem: %w", err)
	}

	source, err := iofs.New(subFS, ".")
	if err != nil {
		return fmt.Errorf("migration source error: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "pgx", driver)
	if err != nil {
		return fmt.Errorf("migrate init error: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migration execution error: %w", err)
	}

	slog.Info("Database migrations completed successfully")
	return nil
}
