package resend

import (
	"fmt"
	"reminder-app/lib/mail"

	resendsdk "github.com/resend/resend-go/v2"
)

var _ mail.Sender = &ResendSender{}

type ResendSender struct {
	ApiKey string
	Domain string
}

func (s *ResendSender) Send(to string, subject string, body string) error {
	client := resendsdk.NewClient(s.ApiKey)
	params := &resendsdk.SendEmailRequest{
		From:    fmt.Sprintf("UchiBot <reminder@%s>", s.Domain),
		To:      []string{to},
		Html:    body,
		Subject: subject,
	}
	_, err := client.Emails.Send(params)
	if err != nil {
		return err
	}

	return nil
}
