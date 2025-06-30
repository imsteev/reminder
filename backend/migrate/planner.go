package migrate

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type Planner struct {
	ID   string
	Up   func(tx *gorm.DB) error
	Down func(tx *gorm.DB) error
}

func NewPlanner(id string, up func(tx *gorm.DB) error, down func(tx *gorm.DB) error) *Planner {
	return &Planner{ID: id, Up: up, Down: down}
}

func (p *Planner) Migration() *gormigrate.Migration {
	return &gormigrate.Migration{ID: p.ID, Migrate: p.Up, Rollback: p.Down}
}
