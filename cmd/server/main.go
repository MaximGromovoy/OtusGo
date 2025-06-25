package main

import (
	server "OtusGo/internal/httpServer"
	"OtusGo/internal/httpServer/handlers"
	"OtusGo/internal/repository/transactionsRepository"
	"OtusGo/internal/service/cbrService"
	"OtusGo/internal/service/exchangeRateCache"
	"OtusGo/internal/service/exchangeService"
	"OtusGo/internal/service/operationsLoggerCache"
	"OtusGo/internal/service/redisCache"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Создаем репозиторий транзакций
	transactionRepo, err := transactionsRepository.NewMongoTransactionRepository()
	if err != nil {
		log.Fatalf("Ошибка создания репозитория транзакций: %v", err)
	}

	// Создаем сервис ЦБ РФ
	cbrSvc := cbrService.NewCBRService()

	ratesCacheConfig := redisCache.CacheConfig{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	}

	logsCacheConfig := redisCache.CacheConfig{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	}

	exchangeRateCache, err := exchangeRateCache.NewExchangeRatesCache(ratesCacheConfig)

	if err != nil {
		log.Fatalf("Ошибка создания кэша курсов валют: %v", err)
	}

	operationsLoggerCache, err := operationsLoggerCache.NewOperationsLoggerCache(logsCacheConfig)

	if err != nil {
		log.Fatalf("Ошибка создания кэша логов операций: %v", err)
	}

	// Создаем сервис обмена валют с комиссией 1%
	exchangeSvc := exchangeService.NewExchangeService(transactionRepo, cbrSvc, exchangeRateCache, operationsLoggerCache, 1.0)

	// Создаем HTTP хендлер
	exchangeHandler := handlers.NewExchangeHandler(exchangeSvc)

	// Создаем HTTP сервер
	httpServer := server.NewHTTPServer(":8080", exchangeHandler)

	// Запускаем сервер в горутине
	go func() {
		log.Println("Запуск HTTP сервера на порту 8080...")
		if err := httpServer.Start(); err != nil {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	// Ожидаем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Завершение работы сервера...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Stop(ctx); err != nil {
		log.Fatalf("Ошибка остановки сервера: %v", err)
	}

	log.Println("Сервер успешно завершен")
}
