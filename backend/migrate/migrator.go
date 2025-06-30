package migrate

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// NewMigrator creates a new gormigrate migrator with all migrations
func NewMigrator(db *gorm.DB) *gormigrate.Gormigrate {

	plans := []*Planner{
		Planner202412291545,
	}

	migrations := make([]*gormigrate.Migration, 0, len(plans))
	for _, plan := range plans {
		migrations = append(migrations, plan.Migration())
	}

	return gormigrate.New(db, gormigrate.DefaultOptions, migrations)
}
