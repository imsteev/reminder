package migrate

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func NewMigrator(db *gorm.DB) *gormigrate.Gormigrate {

	plans := []*MigrationPlan{
		Plan202412291545,
	}

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
