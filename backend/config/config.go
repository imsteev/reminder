package config

import (
	"context"
	"log"

	"github.com/sethvargo/go-envconfig"
)

type ResendConfig struct {
	ApiKey string `env:"RESEND_API_KEY"`
	Domain string `env:"RESEND_DOMAIN"`
}

type Config struct {
	DatabaseURL string `env:"DATABASE_URL,default=postgres://localhost/reminder?sslmode=disable"`
	Port        string `env:"PORT,default=8080"`
	Resend      *ResendConfig
}

func New() *Config {
	var cfg Config
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		log.Fatalf("Failed to process environment variables: %v", err)
	}

	return &cfg
}
