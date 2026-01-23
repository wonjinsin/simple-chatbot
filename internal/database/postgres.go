package database

import (
	"context"
	"database/sql"

	"github.com/wonjinsin/simple-chatbot/internal/config"
	"github.com/wonjinsin/simple-chatbot/pkg/errors"

	// Import pgx driver for PostgreSQL database connectivity
	_ "github.com/jackc/pgx/v5/stdlib"
)

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(cfg *config.Config) (*sql.DB, error) {
	// Open database connection
	db, err := sql.Open("pgx", cfg.GetDatabaseURL())
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to database")
	}

	// Test connection
	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to ping database")
	}

	// Set connection pool settings (optional but recommended)
	db.SetMaxOpenConns(25)                // Maximum number of open connections
	db.SetMaxIdleConns(5)                 // Maximum number of idle connections
	db.SetConnMaxLifetime(0)              // Maximum lifetime of a connection (0 = unlimited)
	db.SetConnMaxIdleTime(0)              // Maximum idle time of a connection (0 = unlimited)

	return db, nil
}
