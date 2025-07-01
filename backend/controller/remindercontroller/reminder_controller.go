package remindercontroller

import (
	"context"
	"errors"
	"fmt"
	"reminder-app/models"
	"reminder-app/workers"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Controller struct {
	db          *gorm.DB
	riverClient *river.Client[pgx.Tx]
}

type Params struct {
	fx.In

	DB    *gorm.DB
	River *river.Client[pgx.Tx]
}

func New(p Params) *Controller {
	return &Controller{db: p.DB, riverClient: p.River}
}

func (rc *Controller) GetReminders(userID int64) ([]models.Reminder, error) {
	var reminders []models.Reminder
	err := rc.db.Where("user_id = ?", userID).Order("start_time").Find(&reminders).Error
	return reminders, err
}

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

func (rc *Controller) CreateReminder(reminder *CreateReminderRequest) (*models.Reminder, error) {
	dbReminder := &models.Reminder{
		UserID:        reminder.UserID,
		Message:       reminder.Message,
		StartTime:     reminder.StartTime,
		Type:          reminder.Type,
		PeriodMinutes: reminder.PeriodMinutes,
		DeliveryType:  reminder.DeliveryType,
	}

	err := rc.db.Create(dbReminder).Error
	if err != nil {
		return nil, err
	}

	isRepeating := reminder.Type == "repeating"

	if isRepeating && reminder.PeriodMinutes <= 0 {
		return nil, errors.New("period minutes must be greater than 0")
	}

	if isRepeating {
		args := workers.ReminderJobArgs{
			ReminderID: int(dbReminder.ID),
		}
		opts := &river.InsertOpts{
			ScheduledAt: dbReminder.StartTime,
			UniqueOpts: river.UniqueOpts{
				ByArgs: true,
			},
		}
		job := river.NewPeriodicJob(
			river.PeriodicInterval(time.Duration(dbReminder.PeriodMinutes)*time.Minute),
			func() (river.JobArgs, *river.InsertOpts) { return args, opts },
			&river.PeriodicJobOpts{
				RunOnStart: true,
			},
		)

		// This adds a handle in memory, not in the database
		handle := rc.riverClient.PeriodicJobs().Add(job)

		fmt.Println("🔄 ADDED PERIODIC JOB", handle)
		dbReminder.JobID = int(handle)

	} else {
		args := workers.ReminderJobArgs{
			ReminderID: int(dbReminder.ID),
		}
		opts := &river.InsertOpts{
			ScheduledAt: dbReminder.StartTime,
		}
		insertResult, err := rc.riverClient.Insert(context.Background(), args, opts)
		if err != nil {
			return nil, err
		}

		dbReminder.JobID = int(insertResult.Job.ID)

	}

	err = rc.db.Save(dbReminder).Error
	if err != nil {
		return nil, err
	}

	return dbReminder, nil
}

func (rc *Controller) UpdateReminder(id int64, reminder *models.Reminder) error {
	err := rc.db.Model(&models.Reminder{}).Where("id = ?", id).Updates(reminder).Error
	return err
}

func (rc *Controller) DeleteReminder(id int64) error {
	var reminder models.Reminder
	if err := rc.db.Where("id = ?", id).First(&reminder).Error; err != nil {
		return err
	}

	if reminder.Type == "repeating" {
		rc.riverClient.PeriodicJobs().Remove(rivertype.PeriodicJobHandle(reminder.JobID))
	} else {
		rc.riverClient.JobDelete(context.Background(), int64(reminder.JobID))
	}

	// Delete the reminder from database
	if err := rc.db.Delete(&reminder).Error; err != nil {
		return err
	}

	return nil
}
