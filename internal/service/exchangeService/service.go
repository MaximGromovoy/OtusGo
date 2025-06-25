package exchangeService

import (
	"OtusGo/internal/interfaces"
	"OtusGo/internal/model/currency"
	models "OtusGo/internal/model/operationLog"
	"OtusGo/internal/transaction"
	"context"
	"fmt"
	"math"
)

type ExchangeService struct {
	transactionRepository interfaces.TransactionRepositoryInterface
	cbrService            interfaces.CBRServiceInterface

	ratesCache      interfaces.ExchangeRatesCacheInterface
	operationsCache interfaces.OperationsLoggerCacheInterface

	commissionRate float64
}

func NewExchangeService(
	transactionRepo interfaces.TransactionRepositoryInterface,
	cbrSvc interfaces.CBRServiceInterface,
	ratesCache interfaces.ExchangeRatesCacheInterface,
	operationsCache interfaces.OperationsLoggerCacheInterface,
	commissionRate float64,
) *ExchangeService {
	return &ExchangeService{
		transactionRepository: transactionRepo,
		cbrService:            cbrSvc,
		ratesCache:            ratesCache,
		operationsCache:       operationsCache,
		commissionRate:        commissionRate,
	}
}

// ExchangeCurrency выполняет обмен валют
func (s *ExchangeService) ExchangeCurrency(ctx context.Context, req *ExchangeRequest) (*ExchangeResult, error) {
	operationLog := models.OperationLog{
		UserID:       req.UserID,
		Operation:    string(transaction.TransactionTypeExchange),
		Status:       string(transaction.TransactionStatusPending),
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
		Amount:       req.FromAmount,
		ErrorMessage: "",
	}

	// Валидируем запрос
	if err := s.validateExchangeRequest(req); err != nil {
		operationLog.Status = string(transaction.TransactionStatusFailed)
		operationLog.ErrorMessage = err.Error()
		s.operationsCache.LogOperation(ctx, &operationLog)

		return nil, fmt.Errorf("invalid exchange request: %w", err)
	}

	// Получаем актуальные курсы валют
	exchangeRate, err := s.getExchangeRate(ctx, req.FromCurrency, req.ToCurrency)
	if err != nil {
		operationLog.Status = string(transaction.TransactionStatusFailed)
		operationLog.ErrorMessage = err.Error()
		s.operationsCache.LogOperation(ctx, &operationLog)

		return nil, fmt.Errorf("failed to get exchange rate: %w", err)
	}

	// Рассчитываем суммы
	toAmountBeforeCommission := req.FromAmount * exchangeRate
	commission := s.calculateCommission(toAmountBeforeCommission)
	toAmountAfterCommission := toAmountBeforeCommission - commission

	// Создаем транзакцию
	tx := transaction.NewExchangeTransaction(
		req.UserID,
		req.FromCurrency,
		req.ToCurrency,
		req.FromAmount,
		toAmountAfterCommission,
		exchangeRate,
		commission,
	)

	operationLog.TransactionID = tx.ID
	operationLog.Rate = exchangeRate
	operationLog.Commission = commission
	operationLog.Timestamp = tx.Timestamp

	// Сохраняем транзакцию
	if err := s.transactionRepository.Add(tx); err != nil {
		operationLog.Status = string(transaction.TransactionStatusFailed)
		operationLog.ErrorMessage = err.Error()

		tx.Fail()
	} else {
		operationLog.Status = string(transaction.TransactionStatusCompleted)

		tx.Complete()
		s.transactionRepository.Update(tx)
	}

	s.operationsCache.LogOperation(ctx, &operationLog)

	result := &ExchangeResult{
		Transaction:    tx,
		ToAmount:       toAmountAfterCommission,
		ExchangeRate:   exchangeRate,
		Commission:     commission,
		CommissionRate: s.commissionRate,
	}

	return result, nil
}

// GetExchangeRate возвращает текущий курс обмена между валютами без выполнения обмена
func (s *ExchangeService) GetExchangeRate(ctx context.Context, fromCurrency, toCurrency string) (float64, error) {
	if !currency.IsCurrencySupported(fromCurrency) {
		return 0, fmt.Errorf("unsupported from currency: %s", fromCurrency)
	}

	if !currency.IsCurrencySupported(toCurrency) {
		return 0, fmt.Errorf("unsupported to currency: %s", toCurrency)
	}

	if fromCurrency == toCurrency {
		return 1.0, nil
	}

	return s.getExchangeRate(ctx, fromCurrency, toCurrency)
}

// CalculateExchangeAmount рассчитывает сумму обмена без создания транзакции
func (s *ExchangeService) CalculateExchangeAmount(ctx context.Context, fromCurrency, toCurrency string, fromAmount float64) (*ExchangeCalculation, error) {
	if fromAmount <= 0 {
		return nil, fmt.Errorf("invalid amount: %f", fromAmount)
	}

	exchangeRate, err := s.GetExchangeRate(ctx, fromCurrency, toCurrency)
	if err != nil {
		return nil, err
	}

	toAmountBeforeCommission := fromAmount * exchangeRate
	commission := s.calculateCommission(toAmountBeforeCommission)
	toAmountAfterCommission := toAmountBeforeCommission - commission

	return &ExchangeCalculation{
		FromCurrency:             fromCurrency,
		ToCurrency:               toCurrency,
		FromAmount:               fromAmount,
		ToAmountBeforeCommission: toAmountBeforeCommission,
		ToAmountAfterCommission:  toAmountAfterCommission,
		ExchangeRate:             exchangeRate,
		Commission:               commission,
		CommissionRate:           s.commissionRate,
	}, nil
}

// GetCachedRates возвращает текущие закэшированные курсы
func (s *ExchangeService) GetCachedRates(ctx context.Context) (map[string]float64, error) {
	return s.ratesCache.GetAll(ctx)
}

// RefreshRates принудительно обновляет курсы валют
func (s *ExchangeService) RefreshRates(ctx context.Context) error {
	// Получаем актуальные курсы из CBR
	currencies, err := s.cbrService.GetCurrencyRates(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch currency rates from CBR: %w", err)
	}

	// Очищаем существующий кэш перед обновлением
	if err := s.ratesCache.Clear(ctx); err != nil {
		return fmt.Errorf("failed to clear rates cache: %w", err)
	}

	// Обновляем кэш новыми курсами
	for i, fromCurrency := range currencies {
		for j, toCurrency := range currencies {
			if i != j {
				fromRate := fromCurrency.GetValue()
				toRate := toCurrency.GetValue()

				if fromRate > 0 && toRate > 0 {
					exchangeRate := fromRate / toRate
					if err := s.ratesCache.Set(ctx, fromCurrency.GetCode(), toCurrency.GetCode(), exchangeRate); err != nil {
						return fmt.Errorf("failed to cache exchange rate for %s to %s: %w",
							fromCurrency.GetCode(), toCurrency.GetCode(), err)
					}
				}
			}
		}
	}

	return nil
}

// validateExchangeRequest валидирует запрос на обмен
func (s *ExchangeService) validateExchangeRequest(req *ExchangeRequest) error {
	if req.UserID <= 0 {
		return fmt.Errorf("invalid user ID: %d", req.UserID)
	}

	if req.FromAmount <= 0 {
		return fmt.Errorf("invalid from amount: %f", req.FromAmount)
	}

	if !currency.IsCurrencySupported(req.FromCurrency) {
		return fmt.Errorf("unsupported from currency: %s", req.FromCurrency)
	}

	if !currency.IsCurrencySupported(req.ToCurrency) {
		return fmt.Errorf("unsupported to currency: %s", req.ToCurrency)
	}

	if req.FromCurrency == req.ToCurrency {
		return fmt.Errorf("cannot exchange same currencies: %s", req.FromCurrency)
	}

	return nil
}

// getExchangeRate получает курс обмена между двумя валютами
func (s *ExchangeService) getExchangeRate(ctx context.Context, fromCurrency, toCurrency string) (float64, error) {

	exist, err := s.ratesCache.Exists(ctx, fromCurrency, toCurrency)
	if err != nil {
		return 0, fmt.Errorf("failed to check if rate exists: %w", err)
	}

	// Если курс уже есть в кэше, возвращаем его
	if exist {
		rate, err := s.ratesCache.Get(ctx, fromCurrency, toCurrency)
		if err != nil {
			return 0, fmt.Errorf("failed to get cached rate: %w", err)
		}
		return rate, nil
	}

	// Если курс не найден в кэше, получаем его из CBR
	currencies, err := s.cbrService.GetCurrencyRates(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch currency rates: %w", err)
	}

	var fromRate, toRate float64

	fmt.Println("-------------------------------------")
	for _, currency := range currencies {
		if currency.GetName() == fromCurrency {
			fromRate = currency.GetValue()
		}
		if currency.GetName() == toCurrency {
			toRate = currency.GetValue()
		}
		if fromRate > 0 && toRate > 0 {
			break
		}
		fmt.Printf("Name: %s, Value: %.4f, Code: %s\n", currency.GetName(), currency.GetValue(), currency.GetCode())
	}
	fmt.Println("-------------------------------------")

	if fromRate == 0 || toRate == 0 {
		return 0, fmt.Errorf("currency rates not found for %s or %s", fromCurrency, toCurrency)
	}

	exchangeRate := fromRate / toRate

	if err := s.ratesCache.Set(ctx, fromCurrency, toCurrency, exchangeRate); err != nil {
		return 0, fmt.Errorf("failed to cache exchange rate: %w", err)
	}

	return exchangeRate, nil
}

// calculateCommission рассчитывает комиссию
func (s *ExchangeService) calculateCommission(amount float64) float64 {
	commission := amount * (s.commissionRate / 100.0)
	// Округляем до 2 знаков после запятой
	return math.Round(commission*100) / 100
}
