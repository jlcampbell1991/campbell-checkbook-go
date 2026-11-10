// Package db provides database initialization and schema management for the Checkbook application.
package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// InitDB initializes the database connection, ensures the required schema exists,
// and returns the *sql.DB instance. It logs connection and schema steps.
func InitDB() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	log.Println("Connected to database")
	if err = ensureSchema(db); err != nil {
		return nil, err
	}
	return db, nil
}

func ensureSchema(db *sql.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS checkbook_accounts (
    name VARCHAR,
    id VARCHAR PRIMARY KEY,
    created_at VARCHAR,
    updated_at VARCHAR
);
CREATE TABLE IF NOT EXISTS checkbook_line_items (
    name VARCHAR,
    balance NUMERIC,
    account_id VARCHAR,
    id VARCHAR PRIMARY KEY,
    active BOOLEAN,
    created_at VARCHAR,
    updated_at VARCHAR
);`
	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to ensure schema: %w", err)
	}
	log.Println("Database schema ensured")
	return nil
}
