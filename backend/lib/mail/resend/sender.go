package resend

import (
	"fmt"
	"reminder-app/lib/mail"
)

var _ mail.Sender = &ResendSender{}

type ResendSender struct {
	ApiKey string
	Domain string
}

func (s *ResendSender) Send(to string, subject string, body string) error {
	fmt.Println("Sending email to", to, "with subject", subject, "and body", body)
	return nil
}
