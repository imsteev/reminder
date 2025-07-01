package models

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type User struct {
	BaseModel `tstype:",extends"`
	Name      string `json:"name" gorm:"not null"`
}

type ContactMethod struct {
	BaseModel   `tstype:",extends"`
	UserID      int64  `json:"user_id" gorm:"not null"`
	Type        string `json:"type" gorm:"not null;type:contact_type"`
	Value       string `json:"value" gorm:"not null"`
	Description string `json:"description"`
}

type Reminder struct {
	BaseModel       `tstype:",extends"`
	UserID          int64     `json:"user_id" gorm:"not null"`
	RiverJobID      int       `json:"river_job_id"`
	ContactMethodID int64     `json:"contact_method_id" gorm:"not null"`
	Body            string    `json:"body" gorm:"not null"`
	StartTime       time.Time `json:"start_time" gorm:"not null"`
	IsRepeating     bool      `json:"is_repeating" gorm:"not null;default:false"`
	PeriodMinutes   int64     `json:"period_minutes" gorm:"not null;default:0"`
}
