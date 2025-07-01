package workers

import (
	"fmt"
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
	fmt.Println("Resend API Key:", p.Config.Resend.ApiKey)
	fmt.Println("Resend Domain:", p.Config.Resend.Domain)
	river.AddWorker(workers, &ReminderJobWorker{
		GormDB: p.DB,
		EmailSender: &resend.ResendSender{
			ApiKey: p.Config.Resend.ApiKey,
			Domain: p.Config.Resend.Domain,
		},
	})
	return workers
}
