package workers

import (
	"fmt"
	"log"
	"reminder-app/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"gorm.io/gorm"
)

func New() *river.Workers {
	workers := river.NewWorkers()
	river.AddWorker(workers, &PeriodicReminderJobWorker{})
	river.AddWorker(workers, &OneTimeReminderJobWorker{})
	return workers
}

func RestorePeriodicJobs(db *gorm.DB, riverClient *river.Client[pgx.Tx]) {
	fmt.Println("Restoring periodic jobs...")

	var reminders []models.Reminder
	err := db.Where("type = ? AND deleted_at IS NULL", "repeating").Find(&reminders).Error
	if err != nil {
		log.Fatalf("Failed to fetch repeating reminders: %v", err)
	}

	for _, reminder := range reminders {
		if reminder.PeriodMinutes <= 0 {
			continue
		}

		period := time.Duration(reminder.PeriodMinutes) * time.Minute

		periodicJob := river.NewPeriodicJob(
			river.PeriodicInterval(period),
			func() (river.JobArgs, *river.InsertOpts) {
				return PeriodicReminderJobArgs{
						ReminderID:  int(reminder.ID),
						PhoneNumber: "", // Will need to get from contact_methods
						Message:     reminder.Message,
					}, &river.InsertOpts{
						ScheduledAt: reminder.StartTime,
					}
			},
			nil,
		)

		handle := riverClient.PeriodicJobs().Add(periodicJob)

		// Update the JobID in the database with the new handle
		err = db.Model(&reminder).Update("job_id", int(handle)).Error
		if err != nil {
			fmt.Printf("Failed to update job_id for reminder %d: %v\n", reminder.ID, err)
		}

		fmt.Printf("Restored periodic job for reminder %d with handle %d\n", reminder.ID, handle)
	}

	fmt.Println("Periodic jobs restored successfully")
}
