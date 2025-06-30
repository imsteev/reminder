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

	if reminder.Type == "repeating" {
		if reminder.PeriodMinutes <= 0 {
			return errors.New("period minutes must be greater than 0")
		}

		period := time.Duration(reminder.PeriodMinutes) * time.Minute

		periodicJob := river.NewPeriodicJob(
			river.PeriodicInterval(period),
			func() (river.JobArgs, *river.InsertOpts) {
				return workers.PeriodicReminderJobArgs{
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
		insertResult, err := rc.riverClient.Insert(context.Background(), workers.PeriodicReminderJobArgs{
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
