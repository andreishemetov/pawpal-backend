package app

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// OpenPostgres opens a pgx-backed *sql.DB, configures the pool, and verifies connectivity.
func OpenPostgres(ctx context.Context, dsn string) (*sql.DB, error) {
	if strings.TrimSpace(dsn) == "" {
		return nil, errors.New("postgres: DATABASE_URL is empty")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
