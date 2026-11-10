package db

import (
    "log"
    "path/filepath"

    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations executes all up migrations located in internal/db/migrations.
// It expects a PostgreSQL DSN as produced by BuildPostgresURL.
func RunMigrations(dsn string) {
    migrationsDir, err := filepath.Abs("internal/db/migrations")
    if err != nil {
        log.Fatal(err)
    }
    sourceURL := "file://" + filepath.ToSlash(migrationsDir)
    log.Println("Running migrations from:", sourceURL)

    m, err := migrate.New(sourceURL, dsn)
    if err != nil {
        log.Fatal("migration init failed:", err)
    }
    if err := m.Up(); err != nil && err.Error() != "no change" {
        log.Fatal("migration failed:", err)
    }
}
