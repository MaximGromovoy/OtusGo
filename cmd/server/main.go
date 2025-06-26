package main

import (
	postgresDatabase "OtusGo/internal/databases/postgres"
	redisDatabase "OtusGo/internal/databases/redis"
	server "OtusGo/internal/httpServer"
	"OtusGo/internal/httpServer/handlers"
	"OtusGo/internal/interfaces"
	"OtusGo/internal/repository/transactionsRepository"
	"OtusGo/internal/service/cbrService"
	"OtusGo/internal/service/currencyService"
	"OtusGo/internal/service/exchangeRateCache"
	exchangeratesprovider "OtusGo/internal/service/exchangeRatesProvider"
	"OtusGo/internal/service/logger"
	transactionsservice "OtusGo/internal/service/transactionsService"
	presets "OtusGo/preset"
	"fmt"
	"log"
)

func main() {
	// Создаем подключения к базам данных
	postgresDB, err := postgresDatabase.NewPostgresDatabase(presets.GetPostgresConfigurationPreset())
	if err != nil {
		fmt.Printf("Ошибка подключения к PostgreSQL: %v\n", err)
		return
	}

	redisDB, err := redisDatabase.NewRedisDatabase(presets.GetRedisConfigurationPreset())
	if err != nil {
		fmt.Printf("Ошибка подключения к Redis: %v\n", err)
		return
	}

	exchangeRatesProvider := createExchangeRatesProvider(redisDB)
	transactionsService := createTransactionsService(postgresDB)
	logger := createLogger(redisDB)

	// Создаем сервис обмена валют
	currencyService := currencyService.NewCurrencyService(exchangeRatesProvider, transactionsService, logger)

	// Создаем HTTP хендлер
	exchangeHandler := handlers.NewExchangeHandler(currencyService, exchangeRatesProvider)

	// Создаем HTTP сервер
	httpServer := server.NewHTTPServer(":8080", exchangeHandler)

	// Запускаем сервер
	log.Println("Запуск HTTP сервера на порту 8080...")
	if err := httpServer.Start(); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

func createExchangeRatesProvider(db *redisDatabase.RedisDatabase) interfaces.ExchangeRatesProviderInterface {
	cbrService := cbrService.NewCBRService()
	exchangeratesCache, err := exchangeRateCache.NewExchangeRatesCache(db)

	if err != nil {
		fmt.Printf("NewExchangeRatesCache error: %v\n", err)

		return nil
	}

	return exchangeratesprovider.NewExchangeRatesProvider(cbrService, exchangeratesCache)
}

func createTransactionsService(db *postgresDatabase.PostgresDB) interfaces.TransactionsServiceInterface {
	transactionsRepository, err := transactionsRepository.NewTransactionsPostgresRepository(db)
	if err != nil {
		fmt.Printf("NewTransactionsPostgresRepository error: %v\n", err)
	}

	return transactionsservice.NewTransactionsService(transactionsRepository)
}

func createLogger(db *redisDatabase.RedisDatabase) interfaces.LoggerInterface {
	logger, err := logger.NewLogger(db)
	if err != nil {
		fmt.Printf("NewOperationsLoggerCache error: %v\n", err)
		return nil
	}

	return logger
}
