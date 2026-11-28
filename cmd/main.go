package main

import (
	"database/sql"
	"log"

	"github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
	_ "github.com/lib/pq"
)

func main() {
	// Load DB configuration from environment (and .env file if present)
	cfg := db.Load()
	dsn := db.BuildPostgresURL(cfg)

	// Test DB connection
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to open DB connection:", err)
	}
	defer conn.Close()
	if err = conn.Ping(); err != nil {
		log.Fatal("failed to ping DB:", err)
	}
	log.Println("DB connection successful")

	// Run migrations (no-op if none present)
	db.RunMigrations(dsn)

	log.Println("application started")
}
