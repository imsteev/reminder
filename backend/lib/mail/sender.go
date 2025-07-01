package mail

import "fmt"

type ResendSender struct {
}

func (s *ResendSender) Send(to string, subject string, body string) error {
	fmt.Println("Sending email to", to, "with subject", subject, "and body", body)
	return nil
}
