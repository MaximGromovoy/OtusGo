package exchange

import (
	"OtusGo/internal/application/interfaces"
	"OtusGo/internal/application/models/transaction"
	"OtusGo/internal/domain/services/exchangeService"
	"context"
	"fmt"
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
		FromCurrency:   fromCurrency,
		ToCurrency:     toCurrency,
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
