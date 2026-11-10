package main

import (
    "log"
    "os"

    "github.com/joho/godotenv"
    "github.com/jlcampbell1991/campbell-checkbook-go/internal/db"
)

func main() {
    // Load environment variables from .env if present.
    _ = godotenv.Load()

    cfg := db.PostgresConfig{
        User:     os.Getenv("POSTGRES_USER"),
        Password: os.Getenv("POSTGRES_PASSWORD"),
        Host:     os.Getenv("POSTGRES_HOST"),
        DB:       os.Getenv("POSTGRES_DB"),
        Port:     os.Getenv("POSTGRES_PORT"),
    }
    dsn := db.BuildPostgresURL(cfg)
    log.Println("Connecting to DB with DSN:", dsn)
    db.RunMigrations(dsn)
    log.Println("Migrations completed")
}
