package riverclient

import (
	"reminder-app/db"
	"reminder-app/jobs"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/samber/do/v2"
)

// NewRiverClient creates River client for background jobs (clean constructor)
func NewRiverClient(dbConnections *db.Connections) (*river.Client[pgx.Tx], error) {
	return New(dbConnections.PGXPool, &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 10},
		},
		Workers: jobs.NewWorkers(),
	})
}

// newRiverClientDI is a wrapper for DI that calls the clean constructor
func newRiverClientDI(i do.Injector) (*river.Client[pgx.Tx], error) {
	dbConnections, err := do.Invoke[*db.Connections](i)
	if err != nil {
		return nil, err
	}
	return NewRiverClient(dbConnections)
}

// Package defines the river client dependency injection package
var Package = do.Package(
	do.Lazy(newRiverClientDI),
)