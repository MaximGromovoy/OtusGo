package currencyService

import (
	"OtusGo/internal/interfaces"
	"OtusGo/internal/model/currency"
	"OtusGo/internal/service/currencyService/operations"
	"OtusGo/internal/service/currencyService/requests"
	"OtusGo/internal/transaction"
	"context"
	"fmt"
	"time"
)

type CurrencyService struct {
	exchangeRatesProvider interfaces.ExchangeRatesProviderInterface
	transactionsService   interfaces.TransactionsServiceInterface
	logger                interfaces.LoggerInterface
}

func NewCurrencyService(exchangeRatesProvider interfaces.ExchangeRatesProviderInterface,
	transactionsService interfaces.TransactionsServiceInterface, logger interfaces.LoggerInterface) *CurrencyService {
	return &CurrencyService{
		exchangeRatesProvider: exchangeRatesProvider,
		transactionsService:   transactionsService,
		logger:                logger,
	}
}

// ExchangeCurrency выполняет обмен валют
func (s *CurrencyService) Exchange(ctx context.Context, req *requests.ExchangeRequest) (*Result, error) {
	start := time.Now()

	userId := req.GetUserId()
	fromCurrency := req.GetFromCurrency()
	toCurrency := req.GetToCurrency()
	fromAmount := req.GetFromAmount()
	commissionRate := req.GetCommissionRate()

	tx := transaction.NewExchangeTransaction(userId, fromCurrency, toCurrency, fromAmount)

	result := &Result{
		Amount:       0,
		Commission:   0,
		ExchangeRate: 0,
		IsSuccessful: true,
	}

	// Валидируем запрос
	if err := validateExchangeRequest(req); err != nil {
		result.IsSuccessful = false
		tx.Status = transaction.TransactionStatusRejected
		s.logger.LogTransaction(ctx, tx, time.Since(start), err.Error())

		return result, fmt.Errorf("request is not valid: %w", err)
	}

	// Получаем курс обмена
	exchangeRate, err := s.exchangeRatesProvider.GetRate(ctx, fromCurrency, toCurrency)
	if err != nil {
		result.IsSuccessful = false
		tx.Status = transaction.TransactionStatusRejected
		s.logger.LogTransaction(ctx, tx, time.Since(start), err.Error())

		return result, fmt.Errorf("failed to get exchange rate: %w", err)
	}

	// Создаем запрос для операции обмена
	exchangeReq := &operations.ExchangeOperationRequest{
		FromAmount:     fromAmount,
		ExchangeRate:   exchangeRate,
		CommissionRate: commissionRate,
	}

	// Выполняем операцию обмена
	operationResult := operations.Exchange(exchangeReq)

	// Заполняем транзакцию результатами обмена
	tx.ToAmount = operationResult.ToAmountAfterCommission
	tx.Commission = operationResult.Commission
	tx.ExchangeRate = exchangeRate

	s.transactionsService.SaveTransaction(tx)

	if tx.Status != transaction.TransactionStatusCompleted {
		result.IsSuccessful = false
	}

	result.Amount = operationResult.ToAmountAfterCommission
	result.Commission = operationResult.Commission
	result.ExchangeRate = exchangeRate

	s.logger.LogTransaction(ctx, tx, time.Since(start), "")

	return result, nil
}

// CalculateExchangeAmount рассчитывает сумму обмена без создания транзакции
func (s *CurrencyService) CalculateExchangeAmount(ctx context.Context, req *requests.CalculateExchangeRequest) (*Result, error) {
	fromCurrency := req.GetFromCurrency()
	toCurrency := req.GetToCurrency()
	fromAmount := req.GetFromAmount()
	commissionRate := req.GetCommissionRate()

	result := &Result{
		Amount:       0,
		Commission:   0,
		ExchangeRate: 0,
		IsSuccessful: true,
	}

	// Валидируем запрос
	if err := validateCalculateExchangeRequest(req); err != nil {
		result.IsSuccessful = false
		return result, fmt.Errorf("request is not valid: %w", err)
	}

	// Получаем курс обмена
	exchangeRate, err := s.exchangeRatesProvider.GetRate(ctx, fromCurrency, toCurrency)
	if err != nil {
		result.IsSuccessful = false
		return result, fmt.Errorf("failed to get exchange rate: %w", err)
	}

	// Создаем запрос для операции обмена
	exchangeReq := &operations.ExchangeOperationRequest{
		FromAmount:     fromAmount,
		ExchangeRate:   exchangeRate,
		CommissionRate: commissionRate,
	}

	// Выполняем операцию обмена
	operationResult := operations.Exchange(exchangeReq)

	result.Amount = operationResult.ToAmountAfterCommission
	result.Commission = operationResult.Commission
	result.ExchangeRate = exchangeRate

	return result, nil
}

func validateUserRequest(r requests.UserRequest) error {
	if r.GetUserId() <= 0 {
		return fmt.Errorf("invalid user ID: %d", r.GetUserId())
	}
	return nil
}

func validateCurrenciesRequest(r requests.CurrenciesRequest) error {
	if !currency.IsCurrencySupported(r.GetFromCurrency()) {
		return fmt.Errorf("unsupported from currency: %s", r.GetFromCurrency())
	}

	if !currency.IsCurrencySupported(r.GetToCurrency()) {
		return fmt.Errorf("unsupported to currency: %s", r.GetToCurrency())
	}

	if r.GetFromCurrency() == r.GetToCurrency() {
		return fmt.Errorf("cannot exchange same currencies: %s - %s", r.GetFromCurrency(), r.GetToCurrency())
	}

	return nil
}

func validateAmountRequest(req requests.AmountRequest) error {
	if req.GetFromAmount() <= 0 {
		return fmt.Errorf("invalid from amount: %f", req.GetFromAmount())
	}
	return nil
}

func validateExchangeRequest(req *requests.ExchangeRequest) error {
	if err := validateUserRequest(req); err != nil {
		return fmt.Errorf("user validation failed: %w", err)
	}

	if err := validateCurrenciesRequest(req); err != nil {
		return fmt.Errorf("currencies validation failed: %w", err)
	}

	if err := validateAmountRequest(req); err != nil {
		return fmt.Errorf("amount validation failed: %w", err)
	}

	return nil
}

func validateCalculateExchangeRequest(req *requests.CalculateExchangeRequest) error {
	if err := validateCurrenciesRequest(req); err != nil {
		return fmt.Errorf("currencies validation failed: %w", err)
	}

	if err := validateAmountRequest(req); err != nil {
		return fmt.Errorf("amount validation failed: %w", err)
	}

	return nil
}
