package controller

import (
	"reminder-app/controller/remindercontroller"
	"reminder-app/db"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/samber/do/v2"
)

// NewAppDI creates the main application controller with dependency injection
func NewAppDI(i do.Injector) (*App, error) {
	dbConnections, err := do.Invoke[*db.Connections](i)
	if err != nil {
		return nil, err
	}

	riverClient, err := do.Invoke[*river.Client[pgx.Tx]](i)
	if err != nil {
		return nil, err
	}

	reminderController, err := do.Invoke[*remindercontroller.Controller](i)
	if err != nil {
		return nil, err
	}

	return NewApp(dbConnections.GORM, riverClient, reminderController), nil
}

// Package defines the app controller dependency injection package
var Package = do.Package(
	do.Lazy(NewAppDI),
)