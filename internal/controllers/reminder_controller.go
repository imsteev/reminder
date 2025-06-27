package controllers

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Reminder struct {
	ID            int        `json:"id" db:"id"`
	UserID        string     `json:"user_id" db:"user_id"`
	Message       string     `json:"message" db:"message"`
	PhoneNumber   string     `json:"phone_number" db:"phone_number"`
	Frequency     int        `json:"frequency" db:"frequency"`
	IntervalHours int        `json:"interval_hours" db:"interval_hours"`
	StartTime     time.Time  `json:"start_time" db:"start_time"`
	EndTime       *time.Time `json:"end_time" db:"end_time"`
	IsActive      bool       `json:"is_active" db:"is_active"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

type ReminderController struct {
	db *pgxpool.Pool
}

func NewReminderController(db *pgxpool.Pool) *ReminderController {
	return &ReminderController{db: db}
}

func (rc *ReminderController) GetReminders(userID string) ([]Reminder, error) {
	rows, err := rc.db.Query(context.Background(),
		"SELECT id, user_id, message, phone_number, frequency, interval_hours, start_time, end_time, is_active, created_at, updated_at FROM reminders WHERE user_id = $1",
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reminders []Reminder
	for rows.Next() {
		var r Reminder
		err := rows.Scan(&r.ID, &r.UserID, &r.Message, &r.PhoneNumber, &r.Frequency, &r.IntervalHours, &r.StartTime, &r.EndTime, &r.IsActive, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			return nil, err
		}
		reminders = append(reminders, r)
	}

	return reminders, nil
}

func (rc *ReminderController) CreateReminder(reminder *Reminder) error {
	reminder.CreatedAt = time.Now()
	reminder.UpdatedAt = time.Now()
	reminder.IsActive = true

	err := rc.db.QueryRow(context.Background(),
		`INSERT INTO reminders (user_id, message, phone_number, frequency, interval_hours, start_time, end_time, is_active, created_at, updated_at) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
		reminder.UserID, reminder.Message, reminder.PhoneNumber, reminder.Frequency, reminder.IntervalHours,
		reminder.StartTime, reminder.EndTime, reminder.IsActive, reminder.CreatedAt, reminder.UpdatedAt,
	).Scan(&reminder.ID)

	return err
}

func (rc *ReminderController) UpdateReminder(id int, reminder *Reminder) error {
	reminder.UpdatedAt = time.Now()

	_, err := rc.db.Exec(context.Background(),
		`UPDATE reminders SET message = $1, phone_number = $2, frequency = $3, interval_hours = $4, 
		 start_time = $5, end_time = $6, is_active = $7, updated_at = $8 WHERE id = $9`,
		reminder.Message, reminder.PhoneNumber, reminder.Frequency, reminder.IntervalHours,
		reminder.StartTime, reminder.EndTime, reminder.IsActive, reminder.UpdatedAt, id,
	)

	return err
}

func (rc *ReminderController) DeleteReminder(id int) error {
	_, err := rc.db.Exec(context.Background(), "DELETE FROM reminders WHERE id = $1", id)
	return err
}
