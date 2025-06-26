package main

import (
	postgresDatabase "OtusGo/internal/databases/postgres"
	redisDatabase "OtusGo/internal/databases/redis"
	"OtusGo/internal/interfaces"
	"OtusGo/internal/repository/transactionsRepository"
	"OtusGo/internal/service/cbrService"
	"OtusGo/internal/service/currencyService"
	"OtusGo/internal/service/currencyService/requests"
	"OtusGo/internal/service/exchangeRateCache"
	exchangeratesprovider "OtusGo/internal/service/exchangeRatesProvider"
	"OtusGo/internal/service/logger"
	transactionsservice "OtusGo/internal/service/transactionsService"
	presets "OtusGo/preset"
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
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

	// Создаем зависимости
	exchangeRatesProvider := createExchangeRatesProvider(redisDB)
	transactionsService := createTransactionsService(postgresDB)
	logger := createLogger(redisDB)

	// Создаем сервис обмена валют
	currencyService := currencyService.NewCurrencyService(exchangeRatesProvider, transactionsService, logger)

	scanner := bufio.NewScanner(os.Stdin)
	ctx := context.Background()

	fmt.Println("=== Консольное приложение обмена валют ===")
	fmt.Println("Доступные команды:")
	fmt.Println("1. exchange - обмен валют")
	fmt.Println("2. calculate - расчет обмена")
	fmt.Println("3. rate - получить курс")
	fmt.Println("4. exit - выход")
	fmt.Println()

	for {
		fmt.Print("Введите команду: ")
		if !scanner.Scan() {
			break
		}

		command := strings.TrimSpace(scanner.Text())
		switch command {
		case "exchange":
			handleExchange(currencyService, scanner, ctx)
		case "calculate":
			handleCalculate(currencyService, scanner, ctx)
		case "rate":
			handleRate(exchangeRatesProvider, scanner, ctx)
		case "exit":
			fmt.Println("Завершение")
			return
		default:
			fmt.Println("Неизвестная команда. Попробуйте снова.")
		}
		fmt.Println()
	}
}

func handleExchange(service *currencyService.CurrencyService, scanner *bufio.Scanner, ctx context.Context) {
	fmt.Print("Введите ID пользователя: ")
	scanner.Scan()
	userID, _ := strconv.Atoi(scanner.Text())

	fmt.Print("Введите валюту исходную (Dollar, Euro, Ruble): ")
	scanner.Scan()
	fromCurrency := strings.TrimSpace(scanner.Text())

	fmt.Print("Введите валюту целевую (Dollar, Euro, Ruble): ")
	scanner.Scan()
	toCurrency := strings.TrimSpace(scanner.Text())

	fmt.Print("Введите сумму для обмена: ")
	scanner.Scan()
	amount, _ := strconv.ParseFloat(scanner.Text(), 64)

	req := requests.NewExchangeRequest(userID, fromCurrency, toCurrency, amount, 0.1)

	result, err := service.Exchange(ctx, req)
	if err != nil {
		fmt.Printf("Ошибка обмена: %v\n", err)
		return
	}

	fmt.Printf("Результат обмена:\n")
	fmt.Printf("  Успешно: %t\n", result.IsSuccessful)
	fmt.Printf("  Сумма к получению: %.2f %s\n", result.Amount, toCurrency)
	fmt.Printf("  Курс обмена: %.4f\n", result.ExchangeRate)
	fmt.Printf("  Комиссия: %.2f\n", result.Commission)
}

func handleCalculate(service *currencyService.CurrencyService, scanner *bufio.Scanner, ctx context.Context) {
	fmt.Print("Введите валюту исходную (Dollar, Euro, Ruble): ")
	scanner.Scan()
	fromCurrency := strings.TrimSpace(scanner.Text())

	fmt.Print("Введите валюту целевую (Dollar, Euro, Ruble): ")
	scanner.Scan()
	toCurrency := strings.TrimSpace(scanner.Text())

	fmt.Print("Введите сумму для расчета: ")
	scanner.Scan()
	amount, _ := strconv.ParseFloat(scanner.Text(), 64)

	req := requests.NewCalculateExchangeRequest(fromCurrency, toCurrency, amount, 0.1)

	result, err := service.CalculateExchangeAmount(ctx, req)
	if err != nil {
		fmt.Printf("Ошибка расчета: %v\n", err)
		return
	}

	fmt.Printf("Результат расчета:\n")
	fmt.Printf("  Сумма к получению: %.2f %s\n", result.Amount, toCurrency)
	fmt.Printf("  Курс обмена: %.4f\n", result.ExchangeRate)
	fmt.Printf("  Комиссия: %.2f\n", result.Commission)
}

func handleRate(provider interfaces.ExchangeRatesProviderInterface, scanner *bufio.Scanner, ctx context.Context) {
	fmt.Print("Введите валюту исходную (Dollar, Euro, Ruble): ")
	scanner.Scan()
	fromCurrency := strings.TrimSpace(scanner.Text())

	fmt.Print("Введите валюту целевую (Dollar, Euro, Ruble): ")
	scanner.Scan()
	toCurrency := strings.TrimSpace(scanner.Text())

	rate, err := provider.GetRate(ctx, fromCurrency, toCurrency)
	if err != nil {
		fmt.Printf("Ошибка получения курса: %v\n", err)
		return
	}

	fmt.Printf("Курс %s/%s: %.4f\n", fromCurrency, toCurrency, rate)
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
