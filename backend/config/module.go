package config

import (
	"os"

	"go.uber.org/fx"
)

// NewConfig creates a new configuration instance (clean constructor)
func NewConfig() *Config {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://localhost/reminder?sslmode=disable"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		DatabaseURL: dbURL,
		Port:        port,
	}
}

var Module = fx.Module("config",
	fx.Provide(NewConfig),
)
