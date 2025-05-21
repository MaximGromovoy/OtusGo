package main

import (
	_ "OtusGo/docs" // Важно импортировать сгенерированную документацию!
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

	storageBaseDir := "." // Текущая директория
	repository, err := currencyRepository.NewCurrencyRepository(storageBaseDir)
	if err != nil {
		panic(fmt.Sprintf("Не удалось инициализировать репозиторий валют: %v", err))
	}

	// Инициализация и запуск HTTP-сервера
	httpServer := server.NewCurrencyServer(repository)
	serverPort := "8080" // Можно получить из конфига или флагов

	go func() {
		log.Printf("Запуск HTTP-сервера на порту %s", serverPort)
		log.Printf("Swagger UI доступен по адресу: http://localhost:%s/swagger/", serverPort)
		if err := httpServer.Start(serverPort); err != nil {
			log.Printf("Ошибка запуска HTTP-сервера: %v", err)
			cancel() // Отменяем контекст, чтобы завершить все горутины
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
