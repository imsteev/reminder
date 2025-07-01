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
	err := w.GormDB.Model(&reminder).Where("id = ?", job.Args.ReminderID).First(&reminder).Error
	if err != nil {
		return fmt.Errorf("failed to get reminder: %w", err)
	}

	var contactMethod models.ContactMethod
	err = w.GormDB.Model(&contactMethod).Where("id = ?", reminder.ContactMethodID).First(&contactMethod).Error
	if err != nil {
		return fmt.Errorf("failed to get contact method: %w", err)
	}

	marshaledReminder, err := json.Marshal(reminder)
	if err != nil {
		return fmt.Errorf("failed to marshal reminder: %w", err)
	}

	switch contactMethod.Type {
	case "email":
		return w.EmailSender.Send(contactMethod.Value, "Reminder", string(marshaledReminder))
	case "phone":
		fmt.Println("Phone number:", contactMethod.Value)
	}

	return nil
}
