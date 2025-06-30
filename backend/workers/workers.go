package workers

import "github.com/riverqueue/river"

func New() *river.Workers {
	workers := river.NewWorkers()
	river.AddWorker(workers, &PeriodicReminderJobWorker{})
	river.AddWorker(workers, &OneTimeReminderJobWorker{})
	return workers
}
