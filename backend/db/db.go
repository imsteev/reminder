package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Connections struct {
	GORM    *gorm.DB
	PGXPool *pgxpool.Pool
}

func New(dbURL string) (*Connections, error) {
	// GORM connection for general database operations
	gormDB, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// PGX Pool connection for River background jobs
	pgxPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		// Get underlying sql.DB from GORM to close it
		if sqlDB, dbErr := gormDB.DB(); dbErr == nil {
			sqlDB.Close()
		}
		return nil, err
	}

	return &Connections{
		GORM:    gormDB,
		PGXPool: pgxPool,
	}, nil
}

func (c *Connections) Close() {
	if c.GORM != nil {
		if sqlDB, err := c.GORM.DB(); err == nil {
			sqlDB.Close()
		}
	}
	if c.PGXPool != nil {
		c.PGXPool.Close()
	}
}

