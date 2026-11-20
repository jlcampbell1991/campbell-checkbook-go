package config

import (
	"log"

	"github.com/joho/godotenv"
)

// Load reads the .env file (if present) and makes the variables available.
// It logs a message if the file is missing but does not treat it as a fatal error.
func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using existing environment variables")
	}
}
