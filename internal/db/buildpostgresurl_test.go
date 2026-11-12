package db

import (
    "net/url"
    "testing"
)

func TestBuildPostgresURL(t *testing.T) {
    cfg := PostgresConfig{
        Host:     "localhost",
        Port:     "5432",
        User:     "testuser",
        Password: "p@ss:word",
        DB:       "testdb",
    }
    got := BuildPostgresURL(cfg)
    expectedPassword := url.QueryEscape(cfg.Password)
    expected := "postgres://" + cfg.User + ":" + expectedPassword + "@" + cfg.Host + ":" + cfg.Port + "/" + cfg.DB + "?sslmode=disable"
    if got != expected {
        t.Fatalf("unexpected postgres URL\n got: %s\n want: %s", got, expected)
    }
}
