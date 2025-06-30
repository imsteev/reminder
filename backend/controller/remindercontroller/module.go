package remindercontroller

import (
	"reminder-app/db"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"go.uber.org/fx"
)

type ControllerParams struct {
	fx.In

	DB    *db.Connections
	River *river.Client[pgx.Tx]
}

// NewController creates a new reminder controller (clean constructor)
func NewController(p ControllerParams) *Controller {
	return NewReminderController(p.DB.GORM, p.River)
}

// Module defines the reminder controller fx module
var Module = fx.Module("remindercontroller",
	fx.Provide(NewController),
)