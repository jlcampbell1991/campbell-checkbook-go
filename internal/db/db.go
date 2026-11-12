package db

import (
	"fmt"
	"net/url"
)

type PostgresConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DB       string
}

// BuildPostgresURL constructs a PostgreSQL connection URL from the given config.
func BuildPostgresURL(cfg PostgresConfig) string {
	escapedPassword := url.QueryEscape(cfg.Password)
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, escapedPassword, cfg.Host, cfg.Port, cfg.DB)
}

// RunMigrations is a placeholder for running DB migrations.
func RunMigrations(dsn string) {
	// No migrations implemented yet.
}
