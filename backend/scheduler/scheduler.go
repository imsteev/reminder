package scheduler

import (
	"context"
	"log"
	"time"

	"reminder-app/controllers"
	"reminder-app/jobs"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
)

type Scheduler struct {
	db    *pgxpool.Pool
	river *river.Client[pgx.Tx]
}

func NewScheduler(db *pgxpool.Pool, riverClient *river.Client[pgx.Tx]) *Scheduler {
	return &Scheduler{
		db:    db,
		river: riverClient,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	log.Println("Scheduler started")

	for {
		select {
		case <-ctx.Done():
			log.Println("Scheduler stopped")
			return
		case <-ticker.C:
			s.processReminders(ctx)
		}
	}
}

func (s *Scheduler) processReminders(ctx context.Context) {
	reminders, err := s.getAllActiveReminders(ctx)
	if err != nil {
		log.Printf("Error fetching reminders: %v", err)
		return
	}

	now := time.Now()

	for _, reminder := range reminders {
		if s.shouldSendReminder(reminder, now) {
			s.scheduleReminderJob(ctx, reminder)
		}
	}
}

func (s *Scheduler) getAllActiveReminders(ctx context.Context) ([]controllers.Reminder, error) {
	rows, err := s.db.Query(ctx,
		"SELECT id, user_id, message, phone_number, frequency, interval_hours, start_time, end_time, is_active, created_at, updated_at FROM reminders WHERE is_active = true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reminders []controllers.Reminder
	for rows.Next() {
		var r controllers.Reminder
		err := rows.Scan(&r.ID, &r.UserID, &r.Message, &r.PhoneNumber, &r.Frequency, &r.IntervalHours, &r.StartTime, &r.EndTime, &r.IsActive, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			return nil, err
		}
		reminders = append(reminders, r)
	}

	return reminders, nil
}

func (s *Scheduler) shouldSendReminder(reminder controllers.Reminder, now time.Time) bool {
	if now.Before(reminder.StartTime) {
		return false
	}

	if reminder.EndTime != nil && now.After(*reminder.EndTime) {
		return false
	}

	timeSinceStart := now.Sub(reminder.StartTime)
	intervalDuration := time.Duration(reminder.IntervalHours) * time.Hour

	if timeSinceStart < intervalDuration {
		return false
	}

	remindersSent := int(timeSinceStart / intervalDuration)
	return remindersSent < reminder.Frequency || reminder.Frequency == -1
}

func (s *Scheduler) scheduleReminderJob(ctx context.Context, reminder controllers.Reminder) {
	jobArgs := jobs.ReminderJobArgs{
		ReminderID:  reminder.ID,
		PhoneNumber: reminder.PhoneNumber,
		Message:     reminder.Message,
	}

	_, err := s.river.Insert(ctx, jobArgs, nil)
	if err != nil {
		log.Printf("Error scheduling reminder job for ID %d: %v", reminder.ID, err)
		return
	}

	log.Printf("Scheduled reminder job for ID %d", reminder.ID)
}
