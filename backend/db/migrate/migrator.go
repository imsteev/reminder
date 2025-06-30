package migrate

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var plans = []*MigrationPlan{
	Plan202412291545,
}

func NewMigrator(db *gorm.DB) *gormigrate.Gormigrate {
	migrations := make([]*gormigrate.Migration, 0, len(plans))
	for _, plan := range plans {
		migrations = append(migrations, makeGormMigration(plan))
	}

	return gormigrate.New(db, gormigrate.DefaultOptions, migrations)
}

func makeGormMigration(plan *MigrationPlan) *gormigrate.Migration {
	return &gormigrate.Migration{
		ID:       plan.ID,
		Migrate:  plan.Up,
		Rollback: plan.Down,
	}
}
