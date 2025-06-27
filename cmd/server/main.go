package main

import (
	"OtusGo/internal/application/services/exchange"
	"OtusGo/internal/domain/services/exchangeService"
	presets "OtusGo/internal/infrastructure"
	transactionsLogger "OtusGo/internal/infrastructure/cache/TransactionsLogger"
	exchangeRatesCache "OtusGo/internal/infrastructure/cache/exchangeRates"
	redisCache "OtusGo/internal/infrastructure/cache/redis"
	"OtusGo/internal/infrastructure/external/cbrClient"
	server "OtusGo/internal/infrastructure/httpServer"
	"OtusGo/internal/infrastructure/httpServer/handlers"
	postgresRepository "OtusGo/internal/infrastructure/repository/postgres"
	transactionsRepository "OtusGo/internal/infrastructure/repository/transactions"
	"fmt"
	"log"
)

func main() {
	// Создаем подключения к базам данных
	postgresDB, err := postgresRepository.NewPostgresDatabase(presets.GetPostgresConfigurationPreset())
	if err != nil {
		fmt.Printf("Ошибка подключения к PostgreSQL: %v\n", err)
		return
	}

	redisDB, err := redisCache.NewRedisDatabase(presets.GetRedisConfigurationPreset())
	if err != nil {
		fmt.Printf("Ошибка подключения к Redis: %v\n", err)
		return
	}

	cbrClient := cbrClient.NewCBRClient()
	exchangeRatesCache, _ := exchangeRatesCache.NewExchangeRatesCache(redisDB)
	exchangeRatesProvider := exchange.NewExchangeRatesProvider(cbrClient, exchangeRatesCache)

	transactionRepository, _ := transactionsRepository.NewTransactionsPostgresRepository(postgresDB)
	exchangeService := exchangeService.NewExchangeService()
	logger, _ := transactionsLogger.NewTransactionsLogger(redisDB)

	exchangeOrchestrator := exchange.NewExchangeOrchestrator(
		exchangeRatesProvider, transactionRepository, exchangeService, logger)

	// Создаем HTTP хендлер
	exchangeHandler := handlers.NewExchangeHandler(exchangeOrchestrator, exchangeRatesProvider)

	// Создаем HTTP сервер
	httpServer := server.NewHTTPServer(":8080", exchangeHandler)

	// Запускаем сервер
	log.Println("Запуск HTTP сервера на порту 8080...")
	if err := httpServer.Start(); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
