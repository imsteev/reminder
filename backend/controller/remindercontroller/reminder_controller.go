package remindercontroller

import (
	"context"
	"errors"
	"fmt"
	"reminder-app/controller/protocol"
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

func (rc *Controller) GetReminders(userID int64, includePast bool) ([]protocol.Reminder, error) {
	var dbReminders []models.Reminder
	query := rc.db.Where("user_id = ?", userID)

	if !includePast {
		// For one-time reminders, exclude past ones. For repeating, always include
		query = query.Where("is_repeating OR (start_time > ?)", time.Now())
	}

	err := query.Order("start_time").Find(&dbReminders).Error
	if err != nil {
		return nil, err
	}

	var protocolReminders []protocol.Reminder
	for _, dbReminder := range dbReminders {
		protocolReminders = append(protocolReminders, protocol.Reminder{
			ID:              int64(dbReminder.ID),
			UserID:          dbReminder.UserID,
			Body:            dbReminder.Body,
			StartTime:       dbReminder.StartTime,
			IsRepeating:     dbReminder.IsRepeating,
			PeriodMinutes:   dbReminder.PeriodMinutes,
			ContactMethodID: dbReminder.ContactMethodID,
		})
	}
	return protocolReminders, err
}

func (rc *Controller) CreateReminder(reminder *protocol.CreateReminderRequest) (*models.Reminder, error) {

	var contactMethod models.ContactMethod
	err := rc.db.Where("user_id = ? and id = ?", reminder.UserID, reminder.ContactMethodID).First(&contactMethod).Error
	if err != nil {
		return nil, errors.New("contact method not found")
	}

	dbReminder := &models.Reminder{
		UserID:          reminder.UserID,
		Body:            reminder.Body,
		StartTime:       reminder.StartTime,
		IsRepeating:     reminder.IsRepeating,
		PeriodMinutes:   reminder.PeriodMinutes,
		ContactMethodID: reminder.ContactMethodID,
	}

	err = rc.db.Create(dbReminder).Error
	if err != nil {
		return nil, err
	}

	if reminder.IsRepeating && reminder.PeriodMinutes <= 0 {
		return nil, errors.New("period minutes must be greater than 0")
	}

	if reminder.IsRepeating {
		args := workers.ReminderJobArgs{
			ReminderID: int(dbReminder.ID),
		}
		opts := &river.InsertOpts{
			ScheduledAt: dbReminder.StartTime,
		}
		job := river.NewPeriodicJob(
			river.PeriodicInterval(time.Duration(dbReminder.PeriodMinutes)*time.Minute),
			func() (river.JobArgs, *river.InsertOpts) { return args, opts },
			nil,
		)

		// This adds a handle in memory, not in the database
		handle := rc.riverClient.PeriodicJobs().Add(job)

		fmt.Println("🔄 ADDED PERIODIC JOB", handle)
		dbReminder.RiverJobID = int(handle)

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

		dbReminder.RiverJobID = int(insertResult.Job.ID)

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

	if reminder.IsRepeating {
		rc.riverClient.PeriodicJobs().Remove(rivertype.PeriodicJobHandle(reminder.RiverJobID))
	} else {
		rc.riverClient.JobDelete(context.Background(), int64(reminder.RiverJobID))
	}

	// Delete the reminder from database
	if err := rc.db.Delete(&reminder).Error; err != nil {
		return err
	}

	return nil
}
