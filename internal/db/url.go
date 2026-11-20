package db

import (
	"fmt"
	"net/url"
)

// BuildPostgresURL constructs a PostgreSQL DSN from the given Config.
// It follows the format: postgres://user:password@host:port/dbname?sslmode=mode
func BuildPostgresURL(cfg Config) string {
	// URL-encode user and password
	userInfo := url.UserPassword(cfg.User, cfg.Password)
	u := &url.URL{
		Scheme: "postgres",
		User:   userInfo,
		Host:   fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Path:   cfg.DB,
	}
	// Add query parameters
	q := u.Query()
	q.Set("sslmode", cfg.SSLMode)
	u.RawQuery = q.Encode()
	return u.String()
}
