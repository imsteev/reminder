package jobs

import (
	"context"
	"fmt"
	"log"

	"github.com/riverqueue/river"
)

type ReminderJobArgs struct {
	ReminderID  int    `json:"reminder_id"`
	PhoneNumber string `json:"phone_number"`
	Message     string `json:"message"`
}

func (ReminderJobArgs) Kind() string { return "reminder" }

type ReminderWorker struct {
	river.WorkerDefaults[ReminderJobArgs]
}

func (w ReminderWorker) Work(ctx context.Context, job *river.Job[ReminderJobArgs]) error {
	log.Printf("Processing reminder job for phone: %s, message: %s", 
		job.Args.PhoneNumber, job.Args.Message)
	
	// TODO: Implement actual SMS sending logic here
	// For now, just log the message that would be sent
	fmt.Printf("📱 SMS to %s: %s\n", job.Args.PhoneNumber, job.Args.Message)
	
	return nil
}