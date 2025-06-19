package main

import (
	transactionRepository "OtusGo/internal/repository/transactionsRepository"
	transactionService "OtusGo/internal/service/transactionService"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var storageBaseDir = "."

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go startSignalHandler(ctx, cancel)

	transactionRepository, err := transactionRepository.NewTransactionRepository(storageBaseDir)

	if err != nil {

	}

	transactionService := transactionService.NewTransactionService(transactionRepository)

	transactionService.Deposit(1, "USD", 100.0)

	<-ctx.Done()
	time.Sleep(1 * time.Second)

	println("App stopped succesful..")
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
