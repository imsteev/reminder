package workers

import (
	"fmt"
	"log"
	"reminder-app/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Params struct {
	fx.In

	DB *gorm.DB
}

func New(p Params) *river.Workers {
	workers := river.NewWorkers()
	river.AddWorker(workers, &ReminderJobWorker{
		GormDB: p.DB,
	})
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
		log.Printf("Creating periodic job for reminder %d with period %v", reminder.ID, period)

		// Capture the reminder in a closure to avoid issues with loop variables
		reminderID := int(reminder.ID)

		periodicJob := river.NewPeriodicJob(
			river.PeriodicInterval(period),
			func() (river.JobArgs, *river.InsertOpts) {
				log.Printf("Periodic job triggered for reminder %d", reminderID)
				return ReminderJobArgs{
						ReminderID: reminderID,
					}, &river.InsertOpts{
						ScheduledAt: reminder.StartTime,
						UniqueOpts: river.UniqueOpts{
							ByArgs: true,
						},
					}
			},
			&river.PeriodicJobOpts{
				RunOnStart: true,
			},
		)

		handle := riverClient.PeriodicJobs().Add(periodicJob)
		log.Printf("Added periodic job for reminder %d, got handle: %d", reminder.ID, handle)

		if handle == 0 {
			log.Printf("WARNING: Got handle 0 for reminder %d - periodic job may not have been added", reminder.ID)
			continue
		}

		// Update the JobID in the database with the new handle
		err = db.Model(&reminder).Update("job_id", int(handle)).Error
		if err != nil {
			log.Printf("Failed to update job_id for reminder %d: %v", reminder.ID, err)
		} else {
			log.Printf("Successfully restored periodic job for reminder %d with handle %d", reminder.ID, handle)
		}
	}

	log.Printf("Successfully restored periodic jobs for %d reminders", len(reminders))
}
