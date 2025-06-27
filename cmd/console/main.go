package main

import (
	"OtusGo/internal/application/services/exchange"
	"OtusGo/internal/domain/services/exchangeService"
	presets "OtusGo/internal/infrastructure"
	exchangeRatesCache "OtusGo/internal/infrastructure/cache/exchangeRates"
	redisCache "OtusGo/internal/infrastructure/cache/redis"
	"OtusGo/internal/infrastructure/cache/transactionsLogger"
	"OtusGo/internal/infrastructure/external/cbrClient"
	postgresRepository "OtusGo/internal/infrastructure/repository/postgres"
	transactionsRepository "OtusGo/internal/infrastructure/repository/transactions"
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	scanner := bufio.NewScanner(os.Stdin)
	ctx := context.Background()

	fmt.Println("=== Консольное приложение обмена валют ===")
	fmt.Println("Доступные команды:")
	fmt.Println("1. exchange - обмен валют")
	fmt.Println("2. rate - получить курс")
	fmt.Println("3. exit - выход")
	fmt.Println()

	for {
		fmt.Print("Введите команду: ")
		if !scanner.Scan() {
			break
		}

		command := strings.TrimSpace(scanner.Text())
		switch command {
		case "exchange":
			handleExchange(exchangeOrchestrator, scanner, ctx)
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

func handleExchange(exchangeOrchestrator *exchange.ExchangeOrchestrator, scanner *bufio.Scanner, ctx context.Context) {
	fmt.Print("Введите ID пользователя: ")
	scanner.Scan()
	userID, _ := strconv.Atoi(scanner.Text())

	fmt.Print("Введите валюту исходную (USD, EUR, RUB): ")
	scanner.Scan()
	fromCurrency := strings.TrimSpace(scanner.Text())

	fmt.Print("Введите валюту целевую (USD, EUR, RUB): ")
	scanner.Scan()
	toCurrency := strings.TrimSpace(scanner.Text())

	fmt.Print("Введите сумму для обмена: ")
	scanner.Scan()
	amount, _ := strconv.ParseFloat(scanner.Text(), 64)

	fmt.Print("Введите размер комиссии: ")
	scanner.Scan()
	commission, _ := strconv.ParseFloat(scanner.Text(), 64)

	exchangeReq := exchange.NewExchangeOrchestratorRequest(
		userID, fromCurrency, toCurrency,
		amount, commission)

	result, err := exchangeOrchestrator.Exchange(ctx, exchangeReq)
	if err != nil {
		fmt.Printf("Ошибка обмена: %v\n", err)
		return
	}

	fmt.Printf("Результат обмена:\n")
	fmt.Printf("  Сумма к получению: %.2f %s\n", result.Amount, toCurrency)
	fmt.Printf("  Курс обмена: %.4f\n", result.ExchangeRate)
	fmt.Printf("  Комиссия: %.2f\n", result.Commission)
}

func handleRate(provider *exchange.ExchangeRatesProvider, scanner *bufio.Scanner, ctx context.Context) {
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
