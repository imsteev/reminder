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
	BaseModel
	Name string `json:"name" gorm:"not null"`
}

type ContactMethod struct {
	BaseModel
	UserID      int64  `json:"user_id" gorm:"not null"`
	Type        string `json:"type" gorm:"not null;type:contact_type"`
	Value       string `json:"value" gorm:"not null"`
	Description string `json:"description"`
}

type Reminder struct {
	BaseModel
	UserID        int64     `json:"user_id" gorm:"not null"`
	Message       string    `json:"message" gorm:"not null"`
	StartTime     time.Time `json:"start_time" gorm:"not null"`
	Type          string    `json:"type" gorm:"not null;type:reminder_type"`
	PeriodMinutes int64     `json:"period_minutes" gorm:"not null;default:0"`
	DeliveryType  string    `json:"delivery_type" gorm:"not null;type:delivery_type"`
}
