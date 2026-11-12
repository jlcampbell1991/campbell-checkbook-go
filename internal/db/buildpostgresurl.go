package db

import (
    "fmt"
    "net/url"
)

// PostgresConfig holds configuration needed to build a PostgreSQL DSN.
type PostgresConfig struct {
    User     string
    Password string
    Host     string
    DB       string
    Port     string
}

// BuildPostgresURL constructs a PostgreSQL connection URL from the given config.
// It mirrors the implementation used in campbell-tube-go.
func BuildPostgresURL(cfg PostgresConfig) string {
    // Default host to localhost if empty.
    host := cfg.Host
    if host == "" {
        host = "localhost"
    }

    user := url.QueryEscape(cfg.User)
    password := url.QueryEscape(cfg.Password)

    // If host starts with '/' treat it as a Unix socket path.
    if len(host) > 0 && host[0] == '/' {
        authority := fmt.Sprintf("%s:%s@/%s", user, password, cfg.DB)
        return fmt.Sprintf("postgres://%s?host=%s&port=%s&sslmode=disable", authority, host, cfg.Port)
    }

    // TCP host.
    authority := fmt.Sprintf("%s:%s@%s:%s/%s", user, password, host, cfg.Port, cfg.DB)
    return "postgres://" + authority + "?sslmode=disable"
}
