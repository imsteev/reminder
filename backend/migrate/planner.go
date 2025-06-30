package migrate

import (
	"gorm.io/gorm"
)

type MigrationPlan struct {
	ID   string
	Up   func(tx *gorm.DB) error
	Down func(tx *gorm.DB) error
}

func NewMigrationPlan(id string, up func(tx *gorm.DB) error, down func(tx *gorm.DB) error) *MigrationPlan {
	return &MigrationPlan{ID: id, Up: up, Down: down}
}
