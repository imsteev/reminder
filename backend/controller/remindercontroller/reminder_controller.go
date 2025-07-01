package remindercontroller

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reminder-app/models"
	"reminder-app/workers"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
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

func (rc *Controller) CreateReminder(reminder *models.Reminder) (*models.Reminder, error) {
	err := rc.db.Create(reminder).Error
	if err != nil {
		return nil, err
	}

	isRepeating := reminder.Type == "repeating"

	if isRepeating && reminder.PeriodMinutes <= 0 {
		return nil, errors.New("period minutes must be greater than 0")
	}

	if isRepeating {
		args := workers.PeriodicReminderJobArgs{
			ReminderID: int(reminder.ID),
		}
		opts := &river.InsertOpts{
			ScheduledAt: reminder.StartTime,
		}
		job := river.NewPeriodicJob(
			river.PeriodicInterval(time.Duration(reminder.PeriodMinutes)*time.Minute),
			func() (river.JobArgs, *river.InsertOpts) { return args, opts },
			nil,
		)

		handle := rc.riverClient.PeriodicJobs().Add(job)

		fmt.Println("🔄 ADDED PERIODIC JOB", handle)
		reminder.JobID = int(handle)

	} else {
		args := workers.OneTimeReminderJobArgs{
			ReminderID: int(reminder.ID),
		}
		opts := &river.InsertOpts{
			ScheduledAt: reminder.StartTime,
		}
		insertResult, err := rc.riverClient.Insert(context.Background(), args, opts)
		if err != nil {
			return nil, err
		}

		reminder.JobID = int(insertResult.Job.ID)

	}

	err = rc.db.Save(reminder).Error
	if err != nil {
		return nil, err
	}

	return reminder, nil
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

	if reminder.JobID > 0 {
		if _, err := rc.riverClient.JobCancel(context.Background(), int64(reminder.JobID)); err != nil {
			return err
		}
	} else {
		log.Println("No job found for reminder", reminder.ID)
	}

	// Delete the reminder from database
	if err := rc.db.Delete(&reminder).Error; err != nil {
		return err
	}

	return nil
}
