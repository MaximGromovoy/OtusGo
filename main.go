package main

import (
	transactionRepository "OtusGo/internal/repository/transactionsRepository"
	"OtusGo/internal/service/cbrService"
	"OtusGo/internal/service/exchangeService"
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	fmt.Println("=== Тестирование сервиса обмена валют ===")

	// Создаем репозиторий транзакций
	transactionRepo, err := transactionRepository.NewMongoTransactionRepository()
	if err != nil {
		log.Fatalf("Ошибка создания репозитория транзакций: %v", err)
	}

	// Создаем сервис ЦБ РФ
	cbrSvc := cbrService.NewCBRService()

	// Создаем сервис обмена валют с комиссией 1%
	exchangeSvc := exchangeService.NewExchangeService(transactionRepo, cbrSvc, nil, nil, 1.0)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Тест 1: Получение текущих курсов валют
	fmt.Println("\n1. Получение текущих курсов валют:")
	testGetExchangeRates(ctx, exchangeSvc)

	// Тест 2: Расчет суммы обмена без создания транзакции
	fmt.Println("\n2. Расчет суммы обмена USD -> EUR:")
	testCalculateExchange(ctx, exchangeSvc)

	// Тест 3: Выполнение обмена валют
	fmt.Println("\n3. Выполнение обмена валют:")
	testExchangeCurrency(ctx, exchangeSvc)

	// Тест 4: Просмотр всех транзакций
	fmt.Println("\n4. Просмотр всех транзакций:")
	testGetAllTransactions(transactionRepo)

	// Тест 5: Тестирование кэша
	fmt.Println("\n5. Тестирование кэша курсов:")
	testCacheRates(ctx, exchangeSvc)

	fmt.Println("\n=== Тестирование завершено ===")
}

func testGetExchangeRates(ctx context.Context, svc *exchangeService.ExchangeService) {
	currencies := []string{"Dollar", "Euro", "Lira", "Ruble"}

	for _, from := range currencies {
		for _, to := range currencies {
			if from != to {
				rate, err := svc.GetExchangeRate(ctx, from, to)
				if err != nil {
					fmt.Printf("Ошибка получения курса %s -> %s: %v\n", from, to, err)
				} else {
					fmt.Printf("Курс %s -> %s: %.4f\n", from, to, rate)
				}
			}
		}
	}
}

func testCalculateExchange(ctx context.Context, svc *exchangeService.ExchangeService) {
	calculation, err := svc.CalculateExchangeAmount(ctx, "Dollar", "Euro", 100)
	if err != nil {
		fmt.Printf("Ошибка расчета: %v\n", err)
		return
	}

	fmt.Printf("Исходная сумма: %.2f %s\n", calculation.FromAmount, calculation.FromCurrency)
	fmt.Printf("Курс обмена: %.4f\n", calculation.ExchangeRate)
	fmt.Printf("Сумма до комиссии: %.2f %s\n", calculation.ToAmountBeforeCommission, calculation.ToCurrency)
	fmt.Printf("Комиссия (%.1f%%): %.2f %s\n", calculation.CommissionRate, calculation.Commission, calculation.ToCurrency)
	fmt.Printf("Итоговая сумма: %.2f %s\n", calculation.ToAmountAfterCommission, calculation.ToCurrency)
}

func testExchangeCurrency(ctx context.Context, svc *exchangeService.ExchangeService) {
	request := &exchangeService.ExchangeRequest{
		UserID:       1,
		FromCurrency: "Dollar",
		ToCurrency:   "Ruble",
		FromAmount:   50,
	}

	result, err := svc.ExchangeCurrency(ctx, request)
	if err != nil {
		fmt.Printf("Ошибка обмена: %v\n", err)
		return
	}

	fmt.Printf("Транзакция ID: %d\n", result.Transaction.ID)
	fmt.Printf("Статус: %s\n", result.Transaction.Status)
	fmt.Printf("Обменяно: %.2f %s -> %.2f %s\n",
		request.FromAmount, request.FromCurrency,
		result.ToAmount, request.ToCurrency)
	fmt.Printf("Курс: %.4f\n", result.ExchangeRate)
	fmt.Printf("Комиссия: %.2f %s\n", result.Commission, request.ToCurrency)
	fmt.Printf("Время: %s\n", result.Transaction.Timestamp.Format("2006-01-02 15:04:05"))

	// Выполним еще один обмен
	request2 := &exchangeService.ExchangeRequest{
		UserID:       2,
		FromCurrency: "Euro",
		ToCurrency:   "Dollar",
		FromAmount:   25,
	}

	result2, err := svc.ExchangeCurrency(ctx, request2)
	if err != nil {
		fmt.Printf("Ошибка второго обмена: %v\n", err)
		return
	}

	fmt.Printf("\nВторая транзакция ID: %d\n", result2.Transaction.ID)
	fmt.Printf("Обменяно: %.2f %s -> %.2f %s\n",
		request2.FromAmount, request2.FromCurrency,
		result2.ToAmount, request2.ToCurrency)
}

func testGetAllTransactions(repo *transactionRepository.TransactionsRepository) {
	transactions := repo.GetAll()

	if len(transactions) == 0 {
		fmt.Println("Транзакций не найдено")
		return
	}

	fmt.Printf("Всего транзакций: %d\n", len(transactions))
	for _, tx := range transactions {
		fmt.Printf("ID: %d, Пользователь: %d, %s -> %s, Статус: %s, Время: %s\n",
			tx.ID, tx.UserID, tx.FromCurrency, tx.ToCurrency, tx.Status,
			tx.Timestamp.Format("2006-01-02 15:04:05"))
	}
}

func testCacheRates(ctx context.Context, svc *exchangeService.ExchangeService) {
	// Принудительно обновляем курсы
	err := svc.RefreshRates(ctx)
	if err != nil {
		fmt.Printf("Ошибка обновления курсов: %v\n", err)
		return
	}

	// Получаем закэшированные курсы
	cachedRates, _ := svc.GetCachedRates(ctx)
	if cachedRates == nil {
		fmt.Println("Кэш пуст или устарел")
		return
	}

	fmt.Println("Закэшированные курсы:")
	for currency, rate := range cachedRates {
		fmt.Printf("%s: %.4f\n", currency, rate)
	}

	// Тестируем скорость работы с кэшем
	start := time.Now()
	_, err = svc.GetExchangeRate(ctx, "Dollar", "Euro")
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("Ошибка получения курса из кэша: %v\n", err)
	} else {
		fmt.Printf("Время получения курса из кэша: %v\n", duration)
	}
}
