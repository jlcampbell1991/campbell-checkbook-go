package db

import (
    "log"
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations runs all up migrations located in internal/db/migrations.
func RunMigrations(dsn string) {
    m, err := migrate.New("file://internal/db/migrations", dsn)
    if err != nil {
        log.Fatalf("migration init error: %v", err)
    }
    if err = m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatalf("migration up error: %v", err)
    }
}
