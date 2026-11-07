package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Port      string
	DBURL     string
	JWTSecret string
}

// Load reads configuration from environment variables (or .env file if loaded elsewhere).
func Load() (*Config, error) {
	// Set defaults
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DB_URL", "")

	// Automatically read environment variables
	viper.AutomaticEnv()

	cfg := &Config{
		Port:      viper.GetString("PORT"),
		DBURL:     viper.GetString("DB_URL"),
		JWTSecret: viper.GetString("JWT_SECRET"),
	}
	if cfg.DBURL == "" {
		return nil, fmt.Errorf("DB_URL environment variable is required")
	}
	return cfg, nil
}
