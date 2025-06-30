package remindercontroller

import (
	"reminder-app/db"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/samber/do/v2"
)

// NewController creates a new reminder controller
func NewController(i do.Injector) (*Controller, error) {
	dbConnections, err := do.Invoke[*db.Connections](i)
	if err != nil {
		return nil, err
	}

	riverClient, err := do.Invoke[*river.Client[pgx.Tx]](i)
	if err != nil {
		return nil, err
	}

	return NewReminderController(dbConnections.GORM, riverClient), nil
}

// Package defines the reminder controller dependency injection package
var Package = do.Package(
	do.Lazy(NewController),
)