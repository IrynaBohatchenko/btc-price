package emailsender

import "context"

type Sender struct {
}

func (s Sender) SendEmails(ctx context.Context) error {
	return nil
}
