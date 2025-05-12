package main

import (
	"OtusGo/internal/model/currency"
	"OtusGo/internal/repository/currencyRepository"
	"OtusGo/internal/service/currencyDistributor"
	"OtusGo/internal/service/currencyRepositoryWatcher"
	"OtusGo/internal/service/randomCurrency"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go startSignalHandler(ctx, cancel)

	storageBaseDir := "." // Текущая директория
	repository, err := currencyRepository.NewCurrencyRepository(storageBaseDir)
	if err != nil {
		panic(fmt.Sprintf("Не удалось инициализировать репозиторий валют: %v", err))
	}

	currencyChannel := make(chan currency.CurrencyInterface)

	go currencyRepositoryWatcher.StartWatching(repository, ctx)
	go randomCurrency.StartRandomCurrencyGenerator(currencyChannel, ctx)
	go currencyDistributor.StartDistributeCurrencyInRepository(currencyChannel, repository, ctx)

	<-ctx.Done()

	close(currencyChannel)

	time.Sleep(1 * time.Second)

	println("Приложение завершено успешно..")
}

func startSignalHandler(ctx context.Context, cancel context.CancelFunc) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case sig := <-signalChannel:
			println("Получен сигнал:", sig)
			cancel()
		case <-ctx.Done():
		}
	}()
}
