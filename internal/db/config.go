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
	SSLMode  string // ssl mode, e.g., disable, require
}

// Load reads environment variables (using the existing config loader to process .env) and returns a Config.
func Load() Config {
	// Ensure .env file is loaded.
	config.Load()

	cfg := Config{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DB:       os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
	if cfg.DB == "" {
		cfg.DB = "campbell_checkbook"
	}
	if cfg.SSLMode == "" {
		cfg.SSLMode = "disable"
	}
	if cfg.Host == "" {
		cfg.Host = "localhost"
	}
	if cfg.Port == "" {
		cfg.Port = "5432"
	}
	if cfg.User == "" {
		log.Fatal("DB_USER must be set")
	}
	return cfg
}
