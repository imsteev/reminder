package workers

import (
	"reminder-app/config"
	"reminder-app/lib/mail/resend"

	"github.com/riverqueue/river"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Params struct {
	fx.In

	DB     *gorm.DB
	Config *config.Config
}

func New(p Params) *river.Workers {
	workers := river.NewWorkers()
	river.AddWorker(workers, &ReminderJobWorker{
		GormDB: p.DB,
		EmailSender: &resend.ResendSender{
			ApiKey: p.Config.Resend.ApiKey,
			Domain: p.Config.Resend.Domain,
		},
	})
	return workers
}
