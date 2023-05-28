package main

import (
	"context"
	"fmt"
	"github.com/btc-price/cmd/price_check/handler"
	"github.com/btc-price/internal/btcpriceservice"
	"github.com/btc-price/internal/coingeckoclient"
	"github.com/btc-price/internal/emailsender"
	"github.com/btc-price/internal/emailstorage"
	"github.com/caarlos0/env/v6"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		panic(fmt.Errorf("read config: %w", err))
	}

	ccx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx, stop := signal.NotifyContext(ccx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	btcPriceSrv := btcpriceservice.NewService(coingeckoclient.NewClient(), emailstorage.NewStorage(), emailsender.NewSender())

	btcPriceHndlr := handler.NewBtcPrice(
		btcPriceSrv,
		logger)
	router := handler.MakeRouter(ctx, btcPriceHndlr)

	httpServer := &http.Server{
		Addr:           cfg.Port,
		Handler:        router,
		ReadTimeout:    cfg.ServerTimeout,
		WriteTimeout:   cfg.ServerTimeout,
		IdleTimeout:    cfg.ServerTimeout,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
	}

	errCh := make(chan error)

	go func() {
		logger.Info("listen and serve", zap.String("address", cfg.Port))

		if err := httpServer.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()

	shutdown := func() {
		stop()
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second) //nolint:gomnd
		defer cancel()

		if err := httpServer.Shutdown(ctxShutdown); err != nil {
			logger.Error("http server: shutdown", zap.Error(err))

			return
		}

		logger.Info("service shutdown: graceful!")
	}

	select {
	case err := <-errCh:
		logger.Error("shutdown catch error", zap.Error(err))
		shutdown()
	case <-ctx.Done():
		logger.Info("shutdown context done")
		shutdown()
	}
}
