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

// Load reads environment variables (after loading .env) and returns a Config.
// It logs a fatal error if the mandatory DB_USER is not set.
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
		log.Fatal("DB_USER environment variable is required")
	}
	// Set sensible defaults if not provided
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
