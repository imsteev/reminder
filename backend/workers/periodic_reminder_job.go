package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"reminder-app/models"

	"github.com/riverqueue/river"
	"gorm.io/gorm"
)

type PeriodicReminderJobArgs struct {
	ReminderID int `json:"reminder_id"`
}

// Kind is the unique string name for this job.
func (PeriodicReminderJobArgs) Kind() string { return "periodic_reminder" }

// PeriodicJobWorker is a job worker for sorting strings.
type PeriodicReminderJobWorker struct {
	river.WorkerDefaults[PeriodicReminderJobArgs]
	GormDB *gorm.DB
}

func (w *PeriodicReminderJobWorker) Work(ctx context.Context, job *river.Job[PeriodicReminderJobArgs]) error {
	fmt.Printf("🔄 EXECUTING PeriodicReminderJobWorker for reminder %d\n", job.Args.ReminderID)

	var reminder models.Reminder
	w.GormDB.Model(&reminder).Where("id = ?", job.Args.ReminderID).First(&reminder)

	if reminder.ID == 0 {
		return fmt.Errorf("reminder not found")
	}

	marshaledReminder, err := json.Marshal(reminder)
	if err != nil {
		return fmt.Errorf("failed to marshal reminder: %w", err)
	}

	fmt.Printf("🔄 MARSHALLED REMINDER: %s\n", string(marshaledReminder))

	fmt.Printf("✅ COMPLETED PeriodicReminderJobWorker for reminder %d\n", job.Args.ReminderID)
	return nil
}
