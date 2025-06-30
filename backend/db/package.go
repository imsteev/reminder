package db

import (
	"reminder-app/config"

	"github.com/samber/do/v2"
)

// NewConnections creates database connections using config (clean constructor)
func NewConnections(cfg *config.Config) (*Connections, error) {
	return New(cfg.DatabaseURL)
}

// newConnectionsDI is a wrapper for DI that calls the clean constructor
func newConnectionsDI(i do.Injector) (*Connections, error) {
	cfg, err := do.Invoke[*config.Config](i)
	if err != nil {
		return nil, err
	}
	return NewConnections(cfg)
}

// Package defines the database dependency injection package
var Package = do.Package(
	do.Lazy(newConnectionsDI),
)