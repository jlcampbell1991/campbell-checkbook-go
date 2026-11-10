package db

import (
    "fmt"
    "net/url"
)

// PostgresConfig holds the configuration needed to build a PostgreSQL DSN.
type PostgresConfig struct {
    User     string
    Password string
    Host     string // hostname or socket path
    DB       string // database name
    Port     string // optional, defaults to 5432
}

// BuildPostgresURL creates a PostgreSQL connection URL from the given config.
// It mirrors the implementation used in the campbell‑tube‑go project.
func BuildPostgresURL(cfg PostgresConfig) string {
    // Default host if empty
    host := cfg.Host
    if host == "" {
        host = "localhost"
    }

    // Escape user and password for URL safety
    user := url.QueryEscape(cfg.User)
    password := url.QueryEscape(cfg.Password)

    // Determine if host is a Unix socket (starts with '/')
    var authority string
    if len(host) > 0 && host[0] == '/' {
        // Unix socket format
        authority = fmt.Sprintf("%s:%s@/%s", user, password, cfg.DB)
        return fmt.Sprintf("postgres://%s?host=%s&port=%s&sslmode=disable", authority, host, cfg.Port)
    }

    // TCP host format
    authority = fmt.Sprintf("%s:%s@%s:%s/%s", user, password, host, cfg.Port, cfg.DB)
    return "postgres://" + authority + "?sslmode=disable"
}
