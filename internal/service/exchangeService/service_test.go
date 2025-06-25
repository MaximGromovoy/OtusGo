package exchangeService

import (
	"OtusGo/internal/mocks"
	"OtusGo/internal/model/currency"
	"OtusGo/internal/transaction"
	"context"
	"testing"
)

// Тесты с использованием моков

var commissionRate = 0.5
var commissionRateFloat = commissionRate / 100.0

// TestExchangeService_ValidateExchangeRequest тестирует все кейсы валидации запроса на обмен валют
// Проверяет:
// - Валидацию UserID (должен быть больше 0)
// - Валидацию FromAmount (должен быть больше 0)
// - Проверку поддержки исходной валюты
// - Проверку поддержки целевой валюты
// - Запрет обмена одинаковых валют
// - Корректное прохождение валидного запроса
func TestExchangeService_ValidateExchangeRequest(t *testing.T) {
	// Создаем моки
	mockRepo := mocks.NewMockTransactionRepository()
	mockCBR := mocks.NewMockCBRService()
	mockExRateCache := mocks.NewMockExchangeRatesCache()
	mockOperationsCache := mocks.NewMockOperationsLoggerCache()

	// Создаем сервис с моками
	exchangeSvc := NewExchangeService(mockRepo, mockCBR, mockExRateCache, mockOperationsCache, commissionRate)

	// Тест кейсы для валидации
	testCases := []struct {
		name        string
		request     *ExchangeRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid request",
			request: &ExchangeRequest{
				UserID:       123,
				FromCurrency: currency.DollarCurrency,
				ToCurrency:   currency.RubleCurrency,
				FromAmount:   100.0,
			},
			expectError: false,
		},
		{
			name: "Invalid UserID - zero",
			request: &ExchangeRequest{
				UserID:       0,
				FromCurrency: currency.DollarCurrency,
				ToCurrency:   currency.RubleCurrency,
				FromAmount:   100.0,
			},
			expectError: true,
			errorMsg:    "invalid user ID: 0",
		},
		{
			name: "Invalid UserID - negative",
			request: &ExchangeRequest{
				UserID:       -1,
				FromCurrency: currency.DollarCurrency,
				ToCurrency:   currency.RubleCurrency,
				FromAmount:   100.0,
			},
			expectError: true,
			errorMsg:    "invalid user ID: -1",
		},
		{
			name: "Invalid FromAmount - zero",
			request: &ExchangeRequest{
				UserID:       123,
				FromCurrency: currency.DollarCurrency,
				ToCurrency:   currency.RubleCurrency,
				FromAmount:   0.0,
			},
			expectError: true,
			errorMsg:    "invalid from amount: 0.000000",
		},
		{
			name: "Invalid FromAmount - negative",
			request: &ExchangeRequest{
				UserID:       123,
				FromCurrency: currency.DollarCurrency,
				ToCurrency:   currency.RubleCurrency,
				FromAmount:   -50.0,
			},
			expectError: true,
			errorMsg:    "invalid from amount: -50.000000",
		},
		{
			name: "Unsupported FromCurrency",
			request: &ExchangeRequest{
				UserID:       123,
				FromCurrency: "Bitcoin",
				ToCurrency:   currency.RubleCurrency,
				FromAmount:   100.0,
			},
			expectError: true,
			errorMsg:    "unsupported from currency: Bitcoin",
		},
		{
			name: "Unsupported ToCurrency",
			request: &ExchangeRequest{
				UserID:       123,
				FromCurrency: currency.DollarCurrency,
				ToCurrency:   "Yen",
				FromAmount:   100.0,
			},
			expectError: true,
			errorMsg:    "unsupported to currency: Yen",
		},
		{
			name: "Same FromCurrency and ToCurrency",
			request: &ExchangeRequest{
				UserID:       123,
				FromCurrency: currency.DollarCurrency,
				ToCurrency:   currency.DollarCurrency,
				FromAmount:   100.0,
			},
			expectError: true,
			errorMsg:    "cannot exchange same currencies: Dollar",
		},
		{
			name: "Empty FromCurrency",
			request: &ExchangeRequest{
				UserID:       123,
				FromCurrency: "",
				ToCurrency:   currency.RubleCurrency,
				FromAmount:   100.0,
			},
			expectError: true,
			errorMsg:    "unsupported from currency: ",
		},
		{
			name: "Empty ToCurrency",
			request: &ExchangeRequest{
				UserID:       123,
				FromCurrency: currency.DollarCurrency,
				ToCurrency:   "",
				FromAmount:   100.0,
			},
			expectError: true,
			errorMsg:    "unsupported to currency: ",
		},
	}

	// Выполняем тесты
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := exchangeSvc.validateExchangeRequest(tc.request)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error for case '%s', but got nil", tc.name)
				} else if err.Error() != tc.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tc.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for case '%s', but got: %v", tc.name, err)
				}
			}
		})
	}
}

// TestExchangeService_ExchangeCurrency_WithMocks тестирует успешный обмен валюты через сервис
// Проверяет:
// - Корректность создания транзакции с правильными параметрами пользователя
// - Правильность расчета суммы до применения комиссии (100 USD * 90 RUB/USD = 9000 RUB)
// - Корректность расчета комиссии (9000 RUB * 0.5% = 45 RUB)
// - Правильность итоговой суммы после вычета комиссии (9000 - 45 = 8955 RUB)
// - Сохранение транзакции в репозитории
func TestExchangeService_ExchangeCurrency_WithMocks(t *testing.T) {
	// Создаем моки
	mockRepo := mocks.NewMockTransactionRepository()
	mockCBR := mocks.NewMockCBRService()
	mockExRateCache := mocks.NewMockExchangeRatesCache()
	mockOperationsCache := mocks.NewMockOperationsLoggerCache()

	ctx := context.Background()

	mockCBR.GetCurrencyRates(ctx)

	// Создаем сервис с моками
	exchangeSvc := NewExchangeService(mockRepo, mockCBR, mockExRateCache, mockOperationsCache, commissionRate)

	// Тестируем обмен USD в RUB
	req := &ExchangeRequest{
		UserID:       123,
		FromCurrency: currency.DollarCurrency,
		ToCurrency:   currency.RubleCurrency,
		FromAmount:   100.0,
	}

	result, err := exchangeSvc.ExchangeCurrency(ctx, req)
	if err != nil {
		t.Fatalf("Exchange failed: %v", err)
	}

	// Проверяем результат
	if result.Transaction == nil {
		t.Fatal("Transaction is nil")
	}

	if result.Transaction.UserID != req.UserID {
		t.Errorf("Expected UserID %d, got %d", req.UserID, result.Transaction.UserID)
	}

	// Ожидаемые значения с курсом 90.0 и комиссией 0.5%
	expectedToAmountBeforeCommission := 100.0 * 90.0                                         // 9000 RUB
	expectedCommission := expectedToAmountBeforeCommission * commissionRateFloat             // 45 RUB
	expectedToAmountAfterCommission := expectedToAmountBeforeCommission - expectedCommission // 8955 RUB

	if result.Commission != expectedCommission {
		t.Errorf("Expected Commission %.2f, got %.2f", expectedCommission, result.Commission)
	}

	if result.ToAmount != expectedToAmountAfterCommission {
		t.Errorf("Expected ToAmount %.2f, got %.2f", expectedToAmountAfterCommission, result.ToAmount)
	}

	if len(mockRepo.GetAll()) != 1 {
		t.Errorf("Expected 1 transaction in repository, got %d", len(mockRepo.GetAll()))
	}
}

// TestExchangeService_ExchangeCurrency_ExpectRepositoryError тестирует поведение сервиса при ошибке репозитория
// Проверяет:
// - Что сервис не возвращает ошибку на уровне API при проблемах с сохранением
// - Что транзакция все равно создается
// - Что статус транзакции устанавливается как "failed" при ошибке сохранения
func TestExchangeService_ExchangeCurrency_ExpectRepositoryError(t *testing.T) {
	// Создаем моки
	mockRepo := mocks.NewMockTransactionRepository()
	mockCBR := mocks.NewMockCBRService()
	mockExRateCache := mocks.NewMockExchangeRatesCache()
	mockOperationsCache := mocks.NewMockOperationsLoggerCache()

	mockRepo.SetShouldAddFail(true) // Устанавливаем флаг для имитации ошибки при добавлении транзакции

	ctx := context.Background()

	mockCBR.GetCurrencyRates(ctx)

	// Создаем сервис с моками
	exchangeSvc := NewExchangeService(mockRepo, mockCBR, mockExRateCache, mockOperationsCache, commissionRate)

	// Тестируем обмен USD в RUB
	req := &ExchangeRequest{
		UserID:       123,
		FromCurrency: currency.DollarCurrency,
		ToCurrency:   currency.RubleCurrency,
		FromAmount:   100.0,
	}

	result, err := exchangeSvc.ExchangeCurrency(ctx, req)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if result.Transaction == nil {
		t.Errorf("Expected transaction to be created, got nil")
	}

	if result.Transaction.Status != transaction.TransactionStatusFailed {
		t.Errorf("Expected transaction status 'failed', got '%s'", result.Transaction.Status)
	}
}

// TestExchangeService_ExchangeCurrency_ExpectedCBRServiceError тестирует обработку ошибок внешнего сервиса курсов валют
// Проверяет:
// - Что при недоступности сервиса ЦБ РФ возвращается ошибка
// - Что результат операции равен nil при критической ошибке
// - Правильность обработки сетевых ошибок и недоступности внешних API
func TestExchangeService_ExchangeCurrency_ExpectedCBRServiceError(t *testing.T) {
	// Создаем моки
	mockRepo := mocks.NewMockTransactionRepository()
	mockCBR := mocks.NewMockCBRService()
	mockExRateCache := mocks.NewMockExchangeRatesCache()
	mockOperationsCache := mocks.NewMockOperationsLoggerCache()

	mockCBR.SetShouldFail(true) // Устанавливаем флаг для имитации ошибки при получении курсов валют

	ctx := context.Background()

	mockCBR.GetCurrencyRates(ctx)

	// Создаем сервис с моками
	exchangeSvc := NewExchangeService(mockRepo, mockCBR, mockExRateCache, mockOperationsCache, commissionRate)

	// Тестируем обмен USD в RUB
	req := &ExchangeRequest{
		UserID:       123,
		FromCurrency: currency.DollarCurrency,
		ToCurrency:   currency.RubleCurrency,
		FromAmount:   100.0,
	}

	result, err := exchangeSvc.ExchangeCurrency(ctx, req)

	if result != nil {
		t.Errorf("Expected nil result on error, got %v", result)
	}

	if err == nil {
		t.Fatalf("Expected ex, but nil: %v", err)
	}
}
