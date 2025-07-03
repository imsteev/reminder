package protocol

import "time"

type CreateReminderRequest struct {
	Body            string    `json:"body"`
	StartTime       time.Time `json:"start_time"`
	IsRepeating     bool      `json:"is_repeating"`
	PeriodMinutes   int64     `json:"period_minutes"`
	ContactMethodID int64     `json:"contact_method_id"`
	PhoneNumber     *string   `json:"phone_number"`
	Email           *string   `json:"email"`
}

type Reminder struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	Body            string    `json:"body"`
	StartTime       time.Time `json:"start_time"`
	IsRepeating     bool      `json:"is_repeating"`
	PeriodMinutes   int64     `json:"period_minutes"`
	ContactMethodID int64     `json:"contact_method_id"`
	PhoneNumber     *string   `json:"phone_number"`
	Email           *string   `json:"email"`
}

type ContactMethod struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"user_id"`
	Type        string `json:"type"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type CreateContactMethodRequest struct {
	UserID      int64  `json:"user_id"`
	Type        string `json:"type"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type UpdateContactMethodRequest struct {
	Type        string `json:"type"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type UpdateReminderRequest struct {
	UserID          int64     `json:"user_id"`
	Body            string    `json:"body"`
	StartTime       time.Time `json:"start_time"`
	IsRepeating     bool      `json:"is_repeating"`
	PeriodMinutes   int64     `json:"period_minutes"`
	ContactMethodID int64     `json:"contact_method_id"`
	PhoneNumber     *string   `json:"phone_number"`
	Email           *string   `json:"email"`
}

type DeleteResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type GetRemindersQuery struct {
	IncludePast bool `json:"include_past" form:"include_past"`
}
