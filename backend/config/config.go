package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/sethvargo/go-envconfig"
)

type ResendConfig struct {
	ApiKey string `env:"RESEND_API_KEY"`
	Domain string `env:"RESEND_DOMAIN"`
}

type Config struct {
	DatabaseURL string `env:"DATABASE_URL"`
	Port        string `env:"PORT,default=8080"`
	Resend      ResendConfig
}

func New() *Config {
	// Debug: Check if DATABASE_URL exists in environment
	fmt.Println("DATABASE_URL from os.Getenv:", os.Getenv("DATABASE_URL"))
	fmt.Println("PORT from os.Getenv:", os.Getenv("PORT"))
	fmt.Println("RESEND_API_KEY from os.Getenv:", os.Getenv("RESEND_API_KEY"))

	var cfg Config
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		log.Fatalf("Failed to process environment variables: %v", err)
	}

	fmt.Printf("Parsed config: %+v\n", cfg)
	fmt.Printf("DatabaseURL: '%s'\n", cfg.DatabaseURL)
	return &cfg
}
