package riverclient

import (
	"fmt"
	"reminder-app/jobs"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"go.uber.org/fx"
)

type Params struct {
	fx.In

	PGXPool *pgxpool.Pool
}

func New(p Params) (*river.Client[pgx.Tx], error) {
	config := &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 10},
		},
		Workers: jobs.NewWorkers(),
	}

	riverClient, err := river.NewClient(riverpgxv5.New(p.PGXPool), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create river client: %w", err)
	}

	return riverClient, nil
}
