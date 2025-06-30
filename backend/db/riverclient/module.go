package riverclient

import (
	"reminder-app/db"
	"reminder-app/jobs"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"go.uber.org/fx"
)

type RiverClientParams struct {
	fx.In

	DB *db.Connections
}

// NewRiverClient creates River client for background jobs (clean constructor)
func NewRiverClient(p RiverClientParams) (*river.Client[pgx.Tx], error) {
	return New(p.DB.PGXPool, &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 10},
		},
		Workers: jobs.NewWorkers(),
	})
}

// Module defines the river client fx module
var Module = fx.Module("riverclient",
	fx.Provide(NewRiverClient),
)