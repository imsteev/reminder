package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"reminder-app/lib/mail"
	"reminder-app/models"

	"github.com/riverqueue/river"
	"gorm.io/gorm"
)

type ReminderJobArgs struct {
	ReminderID int `json:"reminder_id"`
}

func (ReminderJobArgs) Kind() string { return "reminder" }

type EmailSender interface {
	Send(to string, subject string, body string) error
}

type ReminderJobWorker struct {
	river.WorkerDefaults[ReminderJobArgs]
	GormDB      *gorm.DB
	EmailSender mail.Sender
}

func (w *ReminderJobWorker) Work(ctx context.Context, job *river.Job[ReminderJobArgs]) error {
	var reminder models.Reminder
	w.GormDB.Model(&reminder).Where("id = ?", job.Args.ReminderID).First(&reminder)
	if reminder.ID == 0 {
		return fmt.Errorf("reminder not found")
	}

	marshaledReminder, err := json.Marshal(reminder)
	if err != nil {
		return fmt.Errorf("failed to marshal reminder: %w", err)
	}

	w.EmailSender.Send("spchung95@gmail.com", "Reminder", string(marshaledReminder))

	return nil
}
