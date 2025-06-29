package jobs

import (
	"context"
	"fmt"

	"github.com/riverqueue/river"
)

type PeriodicReminderJobArgs struct {
	ReminderID  int    `json:"reminder_id"`
	PhoneNumber string `json:"phone_number"`
	Message     string `json:"message"`
}

// Kind is the unique string name for this job.
func (PeriodicReminderJobArgs) Kind() string { return "periodic_reminder" }

// PeriodicJobWorker is a job worker for sorting strings.
type PeriodicReminderJobWorker struct {
	river.WorkerDefaults[PeriodicReminderJobArgs]
}

func (w *PeriodicReminderJobWorker) Work(ctx context.Context, job *river.Job[PeriodicReminderJobArgs]) error {
	fmt.Printf("This job will run once immediately then approximately once every 15 minutes\n")
	return nil
}
