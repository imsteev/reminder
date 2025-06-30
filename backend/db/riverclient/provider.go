package riverclient

import (
	"reminder-app/db"
	"reminder-app/jobs"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/samber/do/v2"
)

// NewRiverClient creates River client for background jobs
func NewRiverClient(i do.Injector) (*river.Client[pgx.Tx], error) {
	dbConnections, err := do.Invoke[*db.Connections](i)
	if err != nil {
		return nil, err
	}

	return New(dbConnections.PGXPool, &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 10},
		},
		Workers: jobs.NewWorkers(),
	})
}

// Package defines the river client dependency injection package
var Package = do.Package(
	do.Lazy(NewRiverClient),
)