package main

import (
	"OtusGo/internal/model/currency"
	"OtusGo/internal/repository/currencyRepository"
	"OtusGo/internal/service/currencyDistributor"
	"OtusGo/internal/service/currencyRepositoryWatcher"
	"OtusGo/internal/service/randomCurrency"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go startSignalHandler(ctx, cancel)

	repository := currencyRepository.NewCurrencyRepository()
	currencyChannel := make(chan currency.CurrencyInterface)

	go currencyRepositoryWatcher.StartWatching(repository, ctx)
	go randomCurrency.StartRandomCurrencyGenerator(currencyChannel, ctx)
	go currencyDistributor.StartDistributeCurrencyInRepository(currencyChannel, repository, ctx)

	<-ctx.Done()

	close(currencyChannel)

	time.Sleep(1 * time.Second)

	println("Application gracefully stopped.")
}

func startSignalHandler(ctx context.Context, cancel context.CancelFunc) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case sig := <-signalChannel:
			println("Received signal:", sig)
			cancel()
		case <-ctx.Done():
		}
	}()
}
