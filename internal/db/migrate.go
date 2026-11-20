package db

import (
	"log"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations executes all up migrations located in the internal/db/migrations directory.
// The dsn should be a PostgreSQL connection string.
func RunMigrations(dsn string) {
	migrationsPath, err := filepath.Abs(filepath.Join("internal", "db", "migrations"))
	if err != nil {
		log.Fatalf("failed to resolve migrations path: %v", err)
	}
	m, err := migrate.New("file://"+migrationsPath, dsn)
	if err != nil {
		log.Fatalf("migration setup failed: %v", err)
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration failed: %v", err)
	}
	log.Println("migrations applied successfully")
}
