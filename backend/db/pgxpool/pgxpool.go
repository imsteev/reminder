package pgxpool

import (
	"context"
	"reminder-app/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

type Params struct {
	fx.In

	Config *config.Config
}

func New(p Params) (*pgxpool.Pool, error) {
	pgxPool, err := pgxpool.New(context.Background(), p.Config.DatabaseURL)
	if err != nil {
		return nil, err
	}
	return pgxPool, nil
}