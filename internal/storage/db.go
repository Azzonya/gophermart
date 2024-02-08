package storage

import (
	"context"
	"embed"
	"errors"
	"github.com/Azzonya/gophermart/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationsDir embed.FS

func NewPgPool(dsn string) (*pgxpool.Pool, error) {
	err := runMigrations(dsn)
	if err != nil {
		return nil, err
	}

	pgpool, err := pgxpool.New(context.Background(), config.Conf.PgDsn)
	if err != nil {
		return nil, err
	}

	return pgpool, nil
}

func runMigrations(dsn string) error {
	d, err := iofs.New(migrationsDir, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		return err
	}

	if err = m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}

	return nil
}
