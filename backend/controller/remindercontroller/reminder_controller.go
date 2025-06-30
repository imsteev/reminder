package remindercontroller

import (
	"reminder-app/jobs"
	"reminder-app/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"gorm.io/gorm"
)

// Use the Reminder model from models package
type Reminder = models.Reminder

type Controller struct {
	db          *gorm.DB
	riverClient *river.Client[pgx.Tx]
}

func NewReminderController(db *gorm.DB, riverClient *river.Client[pgx.Tx]) *Controller {
	return &Controller{db: db, riverClient: riverClient}
}

func (rc *Controller) GetReminders(userID int64) ([]Reminder, error) {
	var reminders []Reminder
	err := rc.db.Where("user_id = ?", userID).Find(&reminders).Error
	return reminders, err
}

func (rc *Controller) CreateReminder(reminder *Reminder) error {
	// GORM automatically handles CreatedAt and UpdatedAt
	err := rc.db.Create(reminder).Error
	if err != nil {
		return err
	}

	// Only create periodic job for repeating reminders
	if reminder.Type == "repeating" && reminder.PeriodMinutes > 0 {
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

		rc.riverClient.PeriodicJobs().Add(periodicJob)
	}

	return nil
}

func (rc *Controller) UpdateReminder(id int64, reminder *Reminder) error {
	// GORM automatically handles UpdatedAt
	err := rc.db.Model(&Reminder{}).Where("id = ?", id).Updates(reminder).Error
	return err
}

func (rc *Controller) DeleteReminder(id int64) error {
	// GORM soft delete - automatically sets deleted_at
	err := rc.db.Delete(&Reminder{}, id).Error
	return err
}
