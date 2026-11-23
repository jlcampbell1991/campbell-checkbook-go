package db

import "testing"

func TestBuildPostgresURL(t *testing.T) {
    cfg := Config{
        User:     "testuser",
        Password: "p@ss word",
        Host:     "localhost",
        Port:     "5432",
        DB:       "testdb",
        SSLMode:  "disable",
    }
    got := BuildPostgresURL(cfg)
    want := "postgres://testuser:p%40ss%20word@localhost:5432/testdb?sslmode=disable"
    if got != want {
        t.Fatalf("unexpected DSN: got %s, want %s", got, want)
    }
}
