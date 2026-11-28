package config

import (
	"github.com/joho/godotenv"
	"log"
)

// Load reads the .env file (if present) and makes the variables available.
func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using existing environment variables")
	}
}
