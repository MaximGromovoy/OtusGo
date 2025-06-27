package exchange

import (
	"OtusGo/internal/application/interfaces"
	"OtusGo/internal/application/models/transaction"
	"OtusGo/internal/domain/services/exchangeService"
	"context"
	"fmt"
	"strings"
	"time"
)

type ExchangeOrchestrator struct {
	exchangeRatesProvider interfaces.ExchangeRatesProviderInterface
	transactionRepository interfaces.TransactionRepositoryInterface
	exchangeService       *exchangeService.ExchangeService
	logger                interfaces.LoggerInterface
}

func NewExchangeOrchestrator(
	exchangeRatesProvider interfaces.ExchangeRatesProviderInterface,
	transactionRepository interfaces.TransactionRepositoryInterface,
	exchangeService *exchangeService.ExchangeService,
	logger interfaces.LoggerInterface) *ExchangeOrchestrator {
	return &ExchangeOrchestrator{
		exchangeRatesProvider: exchangeRatesProvider,
		transactionRepository: transactionRepository,
		exchangeService:       exchangeService,
		logger:                logger,
	}
}

type ExchangeOrchestratorRequest struct {
	UserId         int
	FromCurrency   string
	ToCurrency     string
	Amount         float64
	CommissionRate float64
}

func NewExchangeOrchestratorRequest(usedId int,
	fromCurrency, toCurrency string,
	amount, commissionRate float64) *ExchangeOrchestratorRequest {
	return &ExchangeOrchestratorRequest{
		UserId:         usedId,
		FromCurrency:   strings.TrimSpace(fromCurrency),
		ToCurrency:     strings.TrimSpace(toCurrency),
		Amount:         amount,
		CommissionRate: commissionRate,
	}
}

type ExchangeOrchestratorResponse struct {
	Amount       float64
	Commission   float64
	ExchangeRate float64
}

func NewExchangeOrchestratorResponse(amount, commission, exchangeRate float64) *ExchangeOrchestratorResponse {
	return &ExchangeOrchestratorResponse{
		Amount:       amount,
		Commission:   commission,
		ExchangeRate: exchangeRate,
	}
}

// Exchange координирует полный процесс обмена валют
func (s *ExchangeOrchestrator) Exchange(ctx context.Context, req *ExchangeOrchestratorRequest) (*ExchangeOrchestratorResponse, error) {
	start := time.Now()

	// Валидация запроса
	if err := s.validateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Создаем транзакцию
	ts := transaction.NewExchangeTransaction(
		req.UserId,
		req.FromCurrency,
		req.ToCurrency,
		req.Amount,
	)

	err := s.transactionRepository.Add(ts)

	if err != nil {
		s.logger.LogTransaction(ctx, ts, time.Since(start), err.Error())
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Получаем курс обмена
	exchangeRate, err := s.exchangeRatesProvider.GetRate(ctx, req.FromCurrency, req.ToCurrency)
	if err != nil {
		ts.Fail()
		s.logger.LogTransaction(ctx, ts, time.Since(start), err.Error())
		s.transactionRepository.Update(ts)
		return nil, fmt.Errorf("failed to get exchange rate: %w", err)
	}

	// Заполняем курс обмена
	ts.ExchangeRate = exchangeRate

	// Выполняем обмен через CurrencyService
	exchangeResponse := s.exchangeService.Exchange(exchangeService.NewExchangeRequest(req.Amount, req.CommissionRate, exchangeRate))

	// Заполняем и обновляем транзакцию
	ts.ToAmount = exchangeResponse.Amount
	ts.Commission = exchangeResponse.Commission
	ts.Complete()
	s.transactionRepository.Update(ts)

	// Логируем успешную операцию
	s.logger.LogTransaction(ctx, ts, time.Since(start), "")

	return NewExchangeOrchestratorResponse(exchangeResponse.Amount,
		exchangeResponse.Commission, exchangeRate), nil
}

// validateRequest выполняет валидацию запроса на обмен валют
func (s *ExchangeOrchestrator) validateRequest(req *ExchangeOrchestratorRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	if req.UserId <= 0 {
		return fmt.Errorf("user ID must be positive, got %d", req.UserId)
	}

	if strings.TrimSpace(req.FromCurrency) == "" {
		return fmt.Errorf("from currency cannot be empty")
	}

	if strings.TrimSpace(req.ToCurrency) == "" {
		return fmt.Errorf("to currency cannot be empty")
	}

	if req.FromCurrency == req.ToCurrency {
		return fmt.Errorf("from currency and to currency cannot be the same: %s", req.FromCurrency)
	}

	if req.Amount < 0 {
		return fmt.Errorf("amount cannot be negative, got %f", req.Amount)
	}

	if req.CommissionRate < 0 {
		return fmt.Errorf("commission rate cannot be negative, got %f", req.CommissionRate)
	}

	if req.CommissionRate >= 100 {
		return fmt.Errorf("commission rate cannot be 100%% or more, got %f", req.CommissionRate)
	}

	// Проверяем формат валютных кодов (должны быть 3 символа)
	if len(strings.TrimSpace(req.FromCurrency)) != 3 {
		return fmt.Errorf("from currency must be 3 characters long, got %s", req.FromCurrency)
	}

	if len(strings.TrimSpace(req.ToCurrency)) != 3 {
		return fmt.Errorf("to currency must be 3 characters long, got %s", req.ToCurrency)
	}

	return nil
}
