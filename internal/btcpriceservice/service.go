package btcpriceservice

import (
	"context"
	"fmt"
	"github.com/btc-price/pkg/btcpricelb"
)

type (
	Service struct {
		coingecko    CoingeckoClient
		emailStorage EmailStorage
		emailSender  EmailSender
	}

	CoingeckoClient interface {
		GetRate(ctx context.Context) (btcpricelb.CoingeckoRate, error)
	}

	EmailStorage interface {
		AddEmail(ctx context.Context, email btcpricelb.Email) error
		ReadOneEmail(ctx context.Context, email btcpricelb.Email) bool
		ReadAllEmails(ctx context.Context) ([]btcpricelb.Email, error)
	}

	EmailSender interface {
		SendEmails(ctx context.Context, emailsList []btcpricelb.Email) error
	}
)

func NewService(coingecko CoingeckoClient, emailStorage EmailStorage, emailSender EmailSender) *Service {
	return &Service{coingecko: coingecko, emailStorage: emailStorage, emailSender: emailSender}
}

func (s *Service) HandleRate(ctx context.Context) (btcpricelb.RateResponse, error) {
	rate, err := s.coingecko.GetRate(ctx)
	if err != nil {
		return btcpricelb.RateResponse{}, err
	}

	return btcpricelb.RateResponse{Rate: float64(rate)}, nil
}

func (s *Service) HandleSubscribe(ctx context.Context, email string) error {
	if ok := s.emailStorage.ReadOneEmail(ctx, btcpricelb.Email(email)); ok {
		return fmt.Errorf(btcpricelb.ErrorEmailExists)
	}

	if err := s.emailStorage.AddEmail(ctx, btcpricelb.Email(email)); err != nil {
		return fmt.Errorf("error adding email %w", err)
	}

	return nil
}

func (s *Service) HandleSendEmails(ctx context.Context) error {
	list, err := s.emailStorage.ReadAllEmails(ctx)
	if err != nil {
		return err
	}
	if err := s.emailSender.SendEmails(ctx, list); err != nil {
		return err
	}

	return nil
}
