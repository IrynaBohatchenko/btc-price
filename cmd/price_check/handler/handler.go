//nolint:lll
package handler

import (
	"context"
	"encoding/json"
	"github.com/btc-price/pkg/btcpricelb"
	"go.uber.org/zap"
	"net/http"
)

type (
	BtcPrice struct {
		srv    BtcPriceService
		logger *zap.Logger
	}

	BtcPriceService interface {
		HandleRate(ctx context.Context) (btcpricelb.RateResponse, error)
		HandleSubscribe(ctx context.Context, email string) error
		HandleSendEmails(ctx context.Context) error
	}
)

func NewBtcPrice(sr BtcPriceService, logger *zap.Logger) *BtcPrice {
	return &BtcPrice{
		srv:    sr,
		logger: logger,
	}
}

func (b *BtcPrice) HandleRate(writer http.ResponseWriter, request *http.Request) {
	resp, err := b.srv.HandleRate(request.Context())
	if err != nil {
		b.logger.Error("error getting rate", zap.Error(err))
		b.write(writer, http.StatusBadRequest, "error getting rate")
		return
	}

	b.write(writer, http.StatusOK, resp.Rate)
}

func (b *BtcPrice) HandleSubscribe(writer http.ResponseWriter, request *http.Request) {
	email := request.FormValue(btcpricelb.EmailForm)

	if err := b.srv.HandleSubscribe(request.Context(), email); err != nil {
		b.logger.Error("error subscribing email", zap.Error(err))
		b.write(writer, http.StatusConflict, "email already exists")
		return
	}

	b.write(writer, http.StatusOK, "E-mail додано")
}

func (b *BtcPrice) HandleSendEmails(writer http.ResponseWriter, request *http.Request) {
	if err := b.srv.HandleSendEmails(request.Context()); err != nil {
		b.logger.Error("error sending emails", zap.Error(err))
		return
	}

	b.write(writer, http.StatusOK, "E-mailʼи відправлено")
}

func (p *BtcPrice) write(writer http.ResponseWriter, statusCode int, data interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	if b, ok := data.([]byte); ok {
		if _, err := writer.Write(b); err != nil {
			p.logger.Error(", write response", zap.Error(err), zap.Any("response", data))

			return
		}
	}

	p.logger.Info("writer.go data=", zap.Any("data", data))
	if err := json.NewEncoder(writer).Encode(data); err != nil {
		p.logger.Error("json encoder, write response", zap.Error(err), zap.Any("response", data))
	}
}
