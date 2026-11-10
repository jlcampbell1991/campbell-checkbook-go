// Package db provides database initialization and migration utilities for the Checkbook application.
package db

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations runs all pending database migrations located in the internal/db/migrations directory.
// It expects a PostgreSQL DSN (Data Source Name) string, for example the value of DATABASE_URL.
// The function logs the source URL and any errors. If there are no pending migrations, it does not treat it as a fatal error.
func RunMigrations(dsn string) {
	// Build the file source URL. The path is relative to the repository root.
	sourceURL := "file://internal/db/migrations"
	log.Println("Running migrations from:", sourceURL)

	m, err := migrate.New(sourceURL, dsn)
	if err != nil {
		log.Fatal(err)
	}
	// Apply all up migrations.
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
	// Success (or no change) – nothing else to do.
}
