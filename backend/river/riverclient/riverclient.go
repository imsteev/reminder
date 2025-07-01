package riverclient

import (
	"context"
	"fmt"
	"reminder-app/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"go.uber.org/fx"
)

type Params struct {
	fx.In

	Config  *config.Config
	Workers *river.Workers
}

func New(p Params) (*river.Client[pgx.Tx], error) {
	pgxPool, err := pgxpool.New(context.Background(), p.Config.DatabaseURL)
	if err != nil {
		return nil, err
	}

	config := &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 10},
		},
		Workers: p.Workers,
	}

	riverClient, err := river.NewClient(riverpgxv5.New(pgxPool), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create river client: %w", err)
	}

	return riverClient, nil
}
