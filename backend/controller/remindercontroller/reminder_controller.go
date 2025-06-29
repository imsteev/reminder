package remindercontroller

import (
	"context"
	"reminder-app/jobs"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
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
	JobID         *int64     `json:"job_id" db:"job_id"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

type Controller struct {
	db          *pgxpool.Pool
	riverClient *river.Client[pgx.Tx]
}

func NewReminderController(db *pgxpool.Pool, riverClient *river.Client[pgx.Tx]) *Controller {
	return &Controller{db: db, riverClient: riverClient}
}

func (rc *Controller) GetReminders(userID string) ([]Reminder, error) {
	rows, err := rc.db.Query(context.Background(),
		"SELECT id, user_id, message, phone_number, frequency, interval_hours, start_time, end_time, is_active, job_id, created_at, updated_at FROM reminders WHERE user_id = $1",
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reminders []Reminder
	for rows.Next() {
		var r Reminder
		err := rows.Scan(&r.ID, &r.UserID, &r.Message, &r.PhoneNumber, &r.Frequency, &r.IntervalHours, &r.StartTime, &r.EndTime, &r.IsActive, &r.JobID, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			return nil, err
		}
		reminders = append(reminders, r)
	}

	return reminders, nil
}

func (rc *Controller) CreateReminder(reminder *Reminder) error {
	reminder.CreatedAt = time.Now()
	reminder.UpdatedAt = time.Now()
	reminder.IsActive = true

	err := rc.db.QueryRow(context.Background(),
		`INSERT INTO reminders (user_id, message, phone_number, frequency, interval_hours, start_time, end_time, is_active, created_at, updated_at) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
		reminder.UserID, reminder.Message, reminder.PhoneNumber, reminder.Frequency, reminder.IntervalHours,
		reminder.StartTime, reminder.EndTime, reminder.IsActive, reminder.CreatedAt, reminder.UpdatedAt,
	).Scan(&reminder.ID)
	if err != nil {
		return err
	}

	// this needs to be improved
	period := time.Duration(reminder.IntervalHours) * time.Hour

	periodicJob := river.NewPeriodicJob(
		river.PeriodicInterval(period),
		func() (river.JobArgs, *river.InsertOpts) {
			return jobs.PeriodicReminderJobArgs{
					ReminderID:  reminder.ID,
					PhoneNumber: reminder.PhoneNumber,
					Message:     reminder.Message,
				}, &river.InsertOpts{
					ScheduledAt: reminder.StartTime,
				}
		},
		nil,
	)

	jobID := rc.riverClient.PeriodicJobs().Add(periodicJob)
	jobIDInt := int64(jobID)
	reminder.JobID = &jobIDInt

	_, err = rc.db.Exec(context.Background(),
		"UPDATE reminders SET job_id = $1 WHERE id = $2",
		reminder.JobID, reminder.ID)

	return err
}

func (rc *Controller) UpdateReminder(id int, reminder *Reminder) error {
	reminder.UpdatedAt = time.Now()

	var oldJobID *int64
	err := rc.db.QueryRow(context.Background(), "SELECT job_id FROM reminders WHERE id = $1", id).Scan(&oldJobID)
	if err != nil {
		return err
	}

	if oldJobID != nil {
		rc.riverClient.PeriodicJobs().Remove(rivertype.PeriodicJobHandle(*oldJobID))
	}

	_, err = rc.db.Exec(context.Background(),
		`UPDATE reminders SET message = $1, phone_number = $2, frequency = $3, interval_hours = $4, 
		 start_time = $5, end_time = $6, is_active = $7, updated_at = $8 WHERE id = $9`,
		reminder.Message, reminder.PhoneNumber, reminder.Frequency, reminder.IntervalHours,
		reminder.StartTime, reminder.EndTime, reminder.IsActive, reminder.UpdatedAt, id,
	)
	if err != nil {
		return err
	}

	if reminder.IsActive {
		period := time.Duration(reminder.IntervalHours) * time.Hour

		// TODO: delete old job

		periodicJob := river.NewPeriodicJob(
			river.PeriodicInterval(period),
			func() (river.JobArgs, *river.InsertOpts) {
				return jobs.PeriodicReminderJobArgs{
						ReminderID:  id,
						PhoneNumber: reminder.PhoneNumber,
						Message:     reminder.Message,
					}, &river.InsertOpts{
						ScheduledAt: time.Now(), // TODO: get this from user input
					}
			},
			nil,
		)

		jobID := rc.riverClient.PeriodicJobs().Add(periodicJob)
		jobIDInt := int64(jobID)
		reminder.JobID = &jobIDInt

		_, err = rc.db.Exec(context.Background(),
			"UPDATE reminders SET job_id = $1 WHERE id = $2",
			reminder.JobID, id)
	} else {
		reminder.JobID = nil
		_, err = rc.db.Exec(context.Background(),
			"UPDATE reminders SET job_id = NULL WHERE id = $1", id)
	}

	return err
}

func (rc *Controller) DeleteReminder(id int) error {
	var jobID *int64
	err := rc.db.QueryRow(context.Background(), "SELECT job_id FROM reminders WHERE id = $1", id).Scan(&jobID)
	if err != nil {
		return err
	}

	if jobID != nil {
		rc.riverClient.PeriodicJobs().Remove(rivertype.PeriodicJobHandle(*jobID))
	}

	_, err = rc.db.Exec(context.Background(), "DELETE FROM reminders WHERE id = $1", id)
	return err
}
