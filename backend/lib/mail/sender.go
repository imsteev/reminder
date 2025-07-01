package mail

type Sender interface {
	Send(to string, subject string, body string) error
}
