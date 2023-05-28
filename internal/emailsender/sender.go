package emailsender

import (
	"context"
	"github.com/btc-price/pkg/btcpricelb"
)

type Sender struct {
}

func NewSender() *Sender {
	return &Sender{}
}

func (s *Sender) SendEmails(ctx context.Context, emailsList []btcpricelb.Email) error {
	for _, _ = range emailsList {
		// interaction with third-party email sending service
	}
	return nil
}
