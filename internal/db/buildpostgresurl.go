package db

import (
	"fmt"
	"net/url"
)

func BuildPostgresURL(cfg PostgresConfig) string {
	// Handle unix socket vs tcp host
	host := cfg.Host
	if host == "" {
		host = "localhost"
	}

	user := url.QueryEscape(cfg.User)
	password := url.QueryEscape(cfg.Password)

	// If host looks like a socket path, format differently
	var authority string
	if len(host) > 0 && host[0] == '/' {
		// Unix socket
		authority = fmt.Sprintf("%s:%s@/%s", user, password, cfg.DB)
		return fmt.Sprintf("postgres://%s?host=%s&port=%s&sslmode=disable",
			authority, host, cfg.Port)
	}

	// TCP host
	authority = fmt.Sprintf("%s:%s@%s:%s/%s",
		user, password, host, cfg.Port, cfg.DB)

	return "postgres://" + authority + "?sslmode=disable"
}

type PostgresConfig struct {
	User     string
	Password string
	Host     string
	DB       string
	Port     string
}
