package workers

import (
	"log"
	"reminder-app/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"gorm.io/gorm"
)

func RestorePeriodicJobs(db *gorm.DB, riverClient *river.Client[pgx.Tx]) error {
	var reminders []models.Reminder
	err := db.Where("is_repeating AND deleted_at IS NULL").Find(&reminders).Error
	if err != nil {
		return err
	}

	for _, reminder := range reminders {
		reminderJob := newPeriodicReminderJob(reminder)
		handle := riverClient.PeriodicJobs().Add(reminderJob)
		if err := db.Model(&reminder).Update("job_id", int(handle)).Error; err != nil {
			log.Printf("Failed to update job_id for reminder %d: %v", reminder.ID, err)
			continue
		}
	}

	return nil
}

func newPeriodicReminderJob(reminder models.Reminder) *river.PeriodicJob {
	return river.NewPeriodicJob(
		river.PeriodicInterval(time.Duration(reminder.PeriodMinutes)*time.Minute),
		func() (river.JobArgs, *river.InsertOpts) {
			return ReminderJobArgs{
					ReminderID: int(reminder.ID),
				}, &river.InsertOpts{
					ScheduledAt: reminder.StartTime,
				}
		},
		nil,
	)
}
