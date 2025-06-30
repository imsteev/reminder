package remindercontroller

import (
	"context"
	"errors"
	"fmt"
	"reminder-app/jobs"
	"reminder-app/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"gorm.io/gorm"
)

// Use the Reminder model from models package

type Controller struct {
	db          *gorm.DB
	riverClient *river.Client[pgx.Tx]
}

func NewReminderController(db *gorm.DB, riverClient *river.Client[pgx.Tx]) *Controller {
	return &Controller{db: db, riverClient: riverClient}
}

func (rc *Controller) GetReminders(userID int64) ([]models.Reminder, error) {
	fmt.Println("GetReminders", userID)
	var reminders []models.Reminder
	err := rc.db.Where("user_id = ?", userID).Find(&reminders).Error
	return reminders, err
}

func (rc *Controller) CreateReminder(reminder *models.Reminder) error {
	err := rc.db.Create(reminder).Error
	if err != nil {
		return err
	}

	// Only create periodic job for repeating reminders
	if reminder.Type == "repeating" {
		if reminder.PeriodMinutes <= 0 {
			return errors.New("period minutes must be greater than 0")
		}

		period := time.Duration(reminder.PeriodMinutes) * time.Minute

		periodicJob := river.NewPeriodicJob(
			river.PeriodicInterval(period),
			func() (river.JobArgs, *river.InsertOpts) {
				return jobs.PeriodicReminderJobArgs{
						ReminderID:  int(reminder.ID),
						PhoneNumber: "", // Will need to get from contact_methods
						Message:     reminder.Message,
					}, &river.InsertOpts{
						ScheduledAt: reminder.StartTime,
					}
			},
			nil,
		)

		handle := rc.riverClient.PeriodicJobs().Add(periodicJob)

		reminder.JobID = int(handle)
		err = rc.db.Save(reminder).Error
		if err != nil {
			return err
		}

		fmt.Println("handle", handle)

	} else {
		insertResult, err := rc.riverClient.Insert(context.Background(), jobs.PeriodicReminderJobArgs{
			ReminderID:  int(reminder.ID),
			PhoneNumber: "", // Will need to get from contact_methods
			Message:     reminder.Message,
		}, &river.InsertOpts{
			ScheduledAt: reminder.StartTime,
		})
		if err != nil {
			return err
		}

		reminder.JobID = int(insertResult.Job.ID)
		err = rc.db.Save(reminder).Error
		if err != nil {
			return err
		}

	}

	return nil
}

func (rc *Controller) UpdateReminder(id int64, reminder *models.Reminder) error {
	err := rc.db.Model(&models.Reminder{}).Where("id = ?", id).Updates(reminder).Error
	return err
}

func (rc *Controller) DeleteReminder(id int64) error {
	err := rc.db.Delete(&models.Reminder{}, id).Error
	return err
}
