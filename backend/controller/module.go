package controller

import (
	"reminder-app/controller/remindercontroller"
	"reminder-app/db"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"go.uber.org/fx"
)

type AppControllerParams struct {
	fx.In

	DB               *db.Connections
	River            *river.Client[pgx.Tx]
	ReminderController *remindercontroller.Controller
}

// NewAppController creates the main application controller (clean constructor)
func NewAppController(p AppControllerParams) *App {
	return NewApp(p.DB.GORM, p.River, p.ReminderController)
}

// Module defines the app controller fx module
var Module = fx.Module("controller",
	fx.Provide(NewAppController),
)
