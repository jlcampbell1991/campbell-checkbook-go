package db

import (
	"net"
	"net/url"
	"strings"
)

// BuildPostgresURL constructs a PostgreSQL DSN from the given Config.
func BuildPostgresURL(cfg Config) string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Path:   cfg.DB,
	}

	query := url.Values{}
	query.Set("sslmode", cfg.SSLMode)

	if strings.HasPrefix(cfg.Host, "/") {
		// Unix socket connection
		query.Set("host", cfg.Host)
	} else {
		// TCP connection
		u.Host = net.JoinHostPort(cfg.Host, cfg.Port)
	}

	u.RawQuery = query.Encode()

	return u.String()
}
