package db

import (
	"fmt"
	"net/url"
)

// BuildPostgresURL builds a PostgreSQL DSN from the given Config.
func BuildPostgresURL(cfg Config) string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Path:   cfg.DB,
	}
	q := u.Query()
	q.Set("sslmode", cfg.SSLMode)
	u.RawQuery = q.Encode()
	return u.String()
}
