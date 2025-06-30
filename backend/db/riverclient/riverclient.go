package riverclient

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

func New(pgxPool *pgxpool.Pool, config *river.Config) (*river.Client[pgx.Tx], error) {
	riverClient, err := river.NewClient(riverpgxv5.New(pgxPool), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create river client: %w", err)
	}

	if err := riverClient.Start(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to start river client: %w", err)
	}

	return riverClient, nil
}
