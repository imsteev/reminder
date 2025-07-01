package workers

import (
	"reminder-app/lib/mail"

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
		GormDB:      p.DB,
		EmailSender: &mail.ResendSender{},
	})
	return workers
}
