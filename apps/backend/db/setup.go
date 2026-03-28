package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// DB is the global database connection pool.
var DB *sql.DB

// Setup connects to PostgreSQL and ensures all tables exist.
// It reads the DATABASE_URL env var, falling back to a local default.
func Setup() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/foreningsliv?sslmode=disable"
	}

	var err error
	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Connection pool settings
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	// Verify the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := DB.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to PostgreSQL")

	// Run migrations
	if err := migrate(ctx); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// Close cleanly shuts down the database connection pool.
func Close() {
	if DB != nil {
		DB.Close()
	}
}

// migrate ensures all required tables exist.
func migrate(ctx context.Context) error {
	migrations := []string{
		// Extensions
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`,

		// Profile table: a person in the system
		`CREATE TABLE IF NOT EXISTS profile (
			id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name       TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		// Auth table: one profile can have multiple auth methods
		// (email/password, Google, Facebook, etc.)
		`CREATE TABLE IF NOT EXISTS auth (
			id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			profile_id UUID NOT NULL REFERENCES profile(id) ON DELETE CASCADE,
			email      TEXT NOT NULL UNIQUE,
			password   TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		// Index for looking up all auth methods for a profile
		`CREATE INDEX IF NOT EXISTS idx_auth_profile_id ON auth(profile_id)`,
	}

	for _, m := range migrations {
		if _, err := DB.ExecContext(ctx, m); err != nil {
			return fmt.Errorf("migration failed: %s: %w", m[:40], err)
		}
	}

	log.Println("Database migrations complete")
	return nil
}
