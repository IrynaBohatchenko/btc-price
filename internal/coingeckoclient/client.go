package coingeckoclient

import (
	"context"
	"encoding/json"
	"github.com/btc-price/pkg/btcpricelb"
	"io"
	"net/http"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) GetRate(ctx context.Context) (btcpricelb.CoingeckoRate, error) {
	resp, err := http.Get(btcpricelb.CoingeckoRatePath)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	answerByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var answer btcpricelb.CoingeckoResponse
	if err := json.Unmarshal(answerByte, &answer); err != nil {
		return 0, err
	}

	return btcpricelb.CoingeckoRate(answer.Bitcoin.Uah), nil
}
