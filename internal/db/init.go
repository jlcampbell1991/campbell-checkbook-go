package db

import (
    "log"
    "os"
)

// InitDB reads environment variables, builds the PostgreSQL DSN, runs migrations, and returns the DSN.
// It is analogous to the setup used in campbell‑tube‑go.
func InitDB() string {
    cfg := PostgresConfig{
        User:     os.Getenv("POSTGRES_USER"),
        Password: os.Getenv("POSTGRES_PASSWORD"),
        Host:     os.Getenv("POSTGRES_HOST"),
        DB:       os.Getenv("POSTGRES_DB"),
        Port:     os.Getenv("POSTGRES_PORT"),
    }
    if cfg.Port == "" {
        cfg.Port = "5432"
    }
    dsn := BuildPostgresURL(cfg)
    log.Println("Connecting to DB with DSN:", dsn)
    RunMigrations(dsn)
    return dsn
}
