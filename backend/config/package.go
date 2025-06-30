package config

import (
	"os"

	"github.com/samber/do/v2"
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

// newConfigDI is a wrapper for DI that calls the clean constructor
func newConfigDI(i do.Injector) (*Config, error) {
	return NewConfig(), nil
}

// Package defines the config dependency injection package
var Package = do.Package(
	do.Lazy(newConfigDI),
)