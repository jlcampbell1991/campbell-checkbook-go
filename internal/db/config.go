package db

import (
    "log"
    "os"
    "github.com/jlcampbell1991/campbell-checkbook-go/internal/config"
)

type Config struct {
    User     string
    Password string
    Host     string
    Port     string
    DB       string // database name
    SSLMode  string // e.g., disable, require
}

// Load reads environment variables (via the config loader) and returns a Config.
func Load() Config {
    // Ensure .env is loaded
    config.Load()
    cfg := Config{
        User:     os.Getenv("DB_USER"),
        Password: os.Getenv("DB_PASSWORD"),
        Host:     os.Getenv("DB_HOST"),
        Port:     os.Getenv("DB_PORT"),
        DB:       os.Getenv("DB_NAME"),
        SSLMode:  os.Getenv("DB_SSLMODE"),
    }
    if cfg.User == "" {
        log.Fatal("DB_USER must be set")
    }
    if cfg.Host == "" {
        cfg.Host = "localhost"
    }
    if cfg.Port == "" {
        cfg.Port = "5432"
    }
    if cfg.DB == "" {
        cfg.DB = "postgres"
    }
    if cfg.SSLMode == "" {
        cfg.SSLMode = "disable"
    }
    return cfg
}
