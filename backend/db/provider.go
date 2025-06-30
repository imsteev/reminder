package db

import (
	"reminder-app/config"

	"github.com/samber/do/v2"
)

// NewConnections creates database connections using config
func NewConnections(i do.Injector) (*Connections, error) {
	cfg, err := do.Invoke[*config.Config](i)
	if err != nil {
		return nil, err
	}
	return New(cfg.DatabaseURL)
}

// Package defines the database dependency injection package
var Package = do.Package(
	do.Lazy(NewConnections),
)