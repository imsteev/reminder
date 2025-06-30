package gorm

import (
	"reminder-app/config"

	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Params struct {
	fx.In

	Config *config.Config
}

func New(p Params) (*gorm.DB, error) {
	gormDB, err := gorm.Open(postgres.Open(p.Config.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	return gormDB, nil
}