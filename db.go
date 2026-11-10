//go:build ignore
// +build ignore

package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Global DB connection
var db *sql.DB

// initDB initializes the database connection and ensures the schema.
func initDB() {
	// Read DATABASE_URL from environment
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatalf("DATABASE_URL environment variable is not set")
	}

	// Open connection using pgx driver
	var err error
	db, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Verify connection
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to database")

	// Ensure schema exists
	ensureSchema()
}

func ensureSchema() {
	// Create table if not exists
	schema := `
CREATE TABLE IF NOT EXISTS checkbook_accounts (
    name VARCHAR,
    id VARCHAR PRIMARY KEY,
    created_at VARCHAR,
    updated_at VARCHAR
);`
	_, err := db.Exec(schema)
	if err != nil {
		log.Fatalf("Failed to ensure schema: %v", err)
	}
	log.Println("Database schema ensured")
}
