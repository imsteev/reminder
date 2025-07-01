package protocol

import "time"

type CreateReminderRequest struct {
	UserID        int64     `json:"user_id"`
	Message       string    `json:"message"`
	StartTime     time.Time `json:"start_time"`
	Type          string    `json:"type"`
	PeriodMinutes int64     `json:"period_minutes"`
	DeliveryType  string    `json:"delivery_type"`
	PhoneNumber   *string   `json:"phone_number"`
	Email         *string   `json:"email"`
}

type Reminder struct {
	ID            int64      `json:"id"`
	UserID        int64      `json:"user_id"`
	Message       string     `json:"message"`
	StartTime     time.Time  `json:"start_time"`
	Type          string     `json:"type"`
	PeriodMinutes int64      `json:"period_minutes"`
	DeliveryType  string     `json:"delivery_type"`
	PhoneNumber   *string    `json:"phone_number"`
	Email         *string    `json:"email"`
	LastRun       *time.Time `json:"last_run"`
	NextRun       *time.Time `json:"next_run"`
}
