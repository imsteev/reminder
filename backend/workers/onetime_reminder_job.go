package workers

import (
	"context"
	"fmt"

	"github.com/riverqueue/river"
)

type OneTimeReminderJobArgs struct {
	ReminderID  int    `json:"reminder_id"`
	PhoneNumber string `json:"phone_number"`
	Message     string `json:"message"`
}

// Kind is the unique string name for this job.
func (OneTimeReminderJobArgs) Kind() string { return "one_time_reminder" }

// OneTimeReminderJobWorker is a job worker for sorting strings.
type OneTimeReminderJobWorker struct {
	river.WorkerDefaults[OneTimeReminderJobArgs]
}

func (w *OneTimeReminderJobWorker) Work(ctx context.Context, job *river.Job[OneTimeReminderJobArgs]) error {
	fmt.Printf("OneTimeReminderJobWorker: %v\n", job.Args)
	return nil
}
