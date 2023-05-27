package emailstorage

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/btc-price/pkg/btcpricelb"
	"log"
	"os"
	"strings"
)

type Storage struct {
}

func NewStorage() *Storage {
	_, err := os.Stat(btcpricelb.StoragePath)
	if err != nil {
		file, err := os.Create(btcpricelb.StoragePath)
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
	return &Storage{}
}

func (s *Storage) AddEmail(ctx context.Context, email btcpricelb.Email) error {
	if s.ReadOneEmail(ctx, email) {
		return fmt.Errorf("email already exists")
	}

	file, err := os.OpenFile(btcpricelb.StoragePath, os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	byteSlice := []byte(fmt.Sprint(email, "\n"))
	_, err = file.Write(byteSlice)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) ReadOneEmail(ctx context.Context, email btcpricelb.Email) bool {
	file, err := os.Open(btcpricelb.StoragePath)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), string(email)) {
			return true
		}
	}

	return false
}

func (s *Storage) ReadAllEmails(ctx context.Context) ([]btcpricelb.Email, error) {
	data, err := os.ReadFile(btcpricelb.StoragePath)
	if err != nil {
		return []btcpricelb.Email{}, err
	}

	var emailsList []btcpricelb.Email
	if err := json.Unmarshal(data, &emailsList); err != nil {
		return []btcpricelb.Email{}, err
	}

	return emailsList, err
}
