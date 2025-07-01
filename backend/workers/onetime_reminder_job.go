package workers

import (
	"context"
	"fmt"

	"github.com/riverqueue/river"
	"gorm.io/gorm"
)

type OneTimeReminderJobArgs struct {
	ReminderID int `json:"reminder_id"`
}

// Kind is the unique string name for this job.
func (OneTimeReminderJobArgs) Kind() string { return "one_time_reminder" }

// OneTimeReminderJobWorker is a job worker for sorting strings.
type OneTimeReminderJobWorker struct {
	river.WorkerDefaults[OneTimeReminderJobArgs]

	GormDB *gorm.DB
}

func (w *OneTimeReminderJobWorker) Work(ctx context.Context, job *river.Job[OneTimeReminderJobArgs]) error {
	fmt.Printf("🔄 EXECUTING OneTimeReminderJobWorker for reminder %d\n", job.Args.ReminderID)

	// TODO: Implement the logic to send the reminder

	fmt.Printf("✅ COMPLETED OneTimeReminderJobWorker for reminder %d\n", job.Args.ReminderID)
	return nil
}
