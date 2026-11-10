package db

import (
	"fmt"
	"net/url"
)

// BuildPostgresURL builds a PostgreSQL DSN from the given configuration.
func BuildPostgresURL(cfg PostgresConfig) string {
	// Default host to localhost if empty.
	host := cfg.Host
	if host == "" {
		host = "localhost"
	}

	user := url.QueryEscape(cfg.User)
	password := url.QueryEscape(cfg.Password)

	// If host looks like a Unix socket path, format accordingly.
	if len(host) > 0 && host[0] == '/' {
		authority := fmt.Sprintf("%s:%s@/%s", user, password, cfg.DB)
		return fmt.Sprintf("postgres://%s?host=%s&port=%s&sslmode=disable", authority, host, cfg.Port)
	}

	// Standard TCP host.
	authority := fmt.Sprintf("%s:%s@%s:%s/%s", user, password, host, cfg.Port, cfg.DB)
	return "postgres://" + authority + "?sslmode=disable"
}

type PostgresConfig struct {
	User     string
	Password string
	Host     string
	DB       string
	Port     string
}
