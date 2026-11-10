package db

import "testing"

func TestBuildPostgresURL_TCP(t *testing.T) {
    cfg := PostgresConfig{User: "user", Password: "pass", Host: "localhost", DB: "mydb", Port: "5432"}
    got := BuildPostgresURL(cfg)
    expected := "postgres://user:pass@localhost:5432/mydb?sslmode=disable"
    if got != expected {
        t.Fatalf("expected %s, got %s", expected, got)
    }
}

func TestBuildPostgresURL_UnixSocket(t *testing.T) {
    cfg := PostgresConfig{User: "user", Password: "pass", Host: "/var/run/postgresql", DB: "mydb", Port: "5432"}
    got := BuildPostgresURL(cfg)
    expected := "postgres://user:pass@/mydb?host=/var/run/postgresql&port=5432&sslmode=disable"
    if got != expected {
        t.Fatalf("expected %s, got %s", expected, got)
    }
}
