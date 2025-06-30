package migrate

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type Planner interface {
	CreateMigration() *gormigrate.Migration
}

// NewMigrator creates a new gormigrate migrator with all migrations
func NewMigrator(db *gorm.DB) *gormigrate.Gormigrate {

	plans := []Planner{
		NewMigrationPlan202412291545(),
	}

	migrations := make([]*gormigrate.Migration, len(plans))
	for i, plan := range plans {
		migrations[i] = plan.CreateMigration()
	}

	return gormigrate.New(db, gormigrate.DefaultOptions, migrations)
}
