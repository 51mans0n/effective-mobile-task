// Package config handles reading app configuration from environment variables.
package config

import (
	"github.com/joho/godotenv"
	"os"
)

// Config holds app configuration loaded from environment variables.
type Config struct {
	AppPort  string
	LogLevel string // debug | info | warn | error
	DBDSN    string
	CacheTTL string
}

// Load loads environment config into a Config struct.
func Load() (*Config, error) {
	_ = godotenv.Load(".env")
	return &Config{
		AppPort:  os.Getenv("APP_PORT"),
		LogLevel: os.Getenv("LOG_LEVEL"),
		DBDSN:    os.Getenv("DB_DSN"),
		CacheTTL: os.Getenv("CACHE_TTL"),
	}, nil
}
