// Package database handles PostgreSQL connection and migrations.
package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/lib/pq"

	"github.com/juanozorio/task-api/internal/config"
)

// Connect establishes a connection pool to PostgreSQL and pings the server.
func Connect(cfg config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	maxOpen, _ := strconv.Atoi(cfg.MaxOpenConns)
	maxIdle, _ := strconv.Atoi(cfg.MaxIdleConns)
	lifetime, _ := time.ParseDuration(cfg.ConnMaxLifetime)

	if maxOpen == 0 {
		maxOpen = 25
	}
	if maxIdle == 0 {
		maxIdle = 25
	}
	if lifetime == 0 {
		lifetime = 5 * time.Minute
	}

	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(lifetime)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}

// RunMigrations executes the initial schema migration.
func RunMigrations(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id           UUID PRIMARY KEY,
		description  TEXT        NOT NULL,
		is_completed BOOLEAN     NOT NULL DEFAULT FALSE,
		created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_tasks_is_completed ON tasks (is_completed);
	`

	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
