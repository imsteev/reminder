package db

import (
	"reminder-app/config"

	"go.uber.org/fx"
)

type ConnectionsParams struct {
	fx.In

	Config *config.Config
}

// NewConnections creates database connections using config (clean constructor)
func NewConnections(p ConnectionsParams) (*Connections, error) {
	return New(p.Config.DatabaseURL)
}

// Module defines the database fx module
var Module = fx.Module("db",
	fx.Provide(NewConnections),
)