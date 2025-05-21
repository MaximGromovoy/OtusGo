package main

import (
	_ "OtusGo/docs"
	"OtusGo/internal/model/currency"
	"OtusGo/internal/repository/currencyRepository"
	"OtusGo/internal/server"
	"OtusGo/internal/service/currencyDistributor"
	"OtusGo/internal/service/currencyRepositoryWatcher"
	"OtusGo/internal/service/randomCurrency"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go startSignalHandler(ctx, cancel)

	storageBaseDir := "."
	repository, err := currencyRepository.NewCurrencyRepository(storageBaseDir)
	if err != nil {
		panic(fmt.Sprintf("Repository init error: %v", err))
	}

	// Инициализация и запуск HTTP-сервера
	httpServer := server.NewCurrencyServer(repository)
	serverPort := "8080"

	go func() {
		log.Printf("Run HTTP-server on port %s", serverPort)
		if err := httpServer.Start(serverPort); err != nil {
			log.Printf("Run server error: %v", err)
			cancel()
		}
	}()

	// Запуск остальных сервисов
	currencyChannel := make(chan currency.CurrencyInterface)

	go currencyRepositoryWatcher.StartWatching(repository, ctx)
	go randomCurrency.StartRandomCurrencyGenerator(currencyChannel, ctx)
	go currencyDistributor.StartDistributeCurrencyInRepository(currencyChannel, repository, ctx)

	<-ctx.Done()

	close(currencyChannel)

	time.Sleep(1 * time.Second)

	println("End of program")
}

func startSignalHandler(ctx context.Context, cancel context.CancelFunc) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case sig := <-signalChannel:
			println("Signal:", sig)
			cancel()
		case <-ctx.Done():
		}
	}()
}
