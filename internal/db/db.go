package db

import (
	"campbell-checkbook/internal/config"
	"database/sql"
	"fmt"
)

// New creates a *sql.DB based on the provided configuration.
func New(cfg *config.Config) (*sql.DB, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}
	// Use pgx driver via stdlib
	db, err := sql.Open("pgx", cfg.DBURL)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return db, nil
}
