package config

import (
	"os"

	"github.com/samber/do/v2"
)

// NewConfig creates a new configuration instance
func NewConfig(i do.Injector) (*Config, error) {
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
	}, nil
}

// Package defines the config dependency injection package
var Package = do.Package(
	do.Lazy(NewConfig),
)