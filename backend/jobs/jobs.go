package jobs

import "github.com/riverqueue/river"

func NewWorkers() *river.Workers {
	workers := river.NewWorkers()
	river.AddWorker(workers, &PeriodicReminderJobWorker{})
	return workers
}
