package exchangeService

import (
	"OtusGo/internal/model/currency"
	"OtusGo/internal/transaction"
	"context"
	"fmt"
	"math"
	"time"
)

// ExchangeCurrency выполняет обмен валют
func (s *ExchangeService) ExchangeCurrency(ctx context.Context, req *ExchangeRequest) (*ExchangeResult, error) {
	// Валидируем запрос
	if err := s.validateExchangeRequest(req); err != nil {
		return nil, fmt.Errorf("invalid exchange request: %w", err)
	}

	// Получаем актуальные курсы валют
	exchangeRate, err := s.getExchangeRate(ctx, req.FromCurrency, req.ToCurrency)
	if err != nil {
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

	// Сохраняем транзакцию
	if err := s.transactionRepository.Add(tx); err != nil {
		tx.Fail()
	} else {
		tx.Complete()
	}

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
func (s *ExchangeService) GetCachedRates() map[string]float64 {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	if !s.isCacheValid() {
		return nil
	}

	// Создаем копию для безопасности
	rates := make(map[string]float64)
	for k, v := range s.ratesCache {
		rates[k] = v
	}

	return rates
}

// RefreshRates принудительно обновляет курсы валют
func (s *ExchangeService) RefreshRates(ctx context.Context) error {
	return s.updateCurrencyRatesFromCBR(ctx)
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
	// Проверяем кэш
	if s.isCacheValid() {
		s.cacheMutex.RLock()
		fromRate, fromExists := s.ratesCache[fromCurrency]
		toRate, toExists := s.ratesCache[toCurrency]
		s.cacheMutex.RUnlock()

		if fromExists && toExists {
			exchangeRate := fromRate / toRate
			return exchangeRate, nil
		}
	}

	// Обновляем курсы с ЦБ РФ
	if err := s.updateCurrencyRatesFromCBR(ctx); err != nil {
		return 0, fmt.Errorf("failed to update currency rates: %w", err)
	}

	// Получаем курсы из кэша
	s.cacheMutex.RLock()
	fromRate, fromExists := s.ratesCache[fromCurrency]
	toRate, toExists := s.ratesCache[toCurrency]
	s.cacheMutex.RUnlock()

	if !fromExists {
		return 0, fmt.Errorf("rate for %s not found", fromCurrency)
	}
	if !toExists {
		return 0, fmt.Errorf("rate for %s not found", toCurrency)
	}

	// Рассчитываем кросс-курс
	exchangeRate := fromRate / toRate

	return exchangeRate, nil
}

// updateCurrencyRatesFromCBR обновляет курсы валют с сайта ЦБ РФ
func (s *ExchangeService) updateCurrencyRatesFromCBR(ctx context.Context) error {
	currencies, err := s.cbrService.GetCurrencyRates(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch currency rates from CBR: %w", err)
	}

	// Обновляем кэш
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	// Очищаем старый кэш
	s.ratesCache = make(map[string]float64)

	// Добавляем новые курсы
	for _, curr := range currencies {
		s.ratesCache[curr.GetName()] = curr.GetValue()
	}

	// Рубль всегда имеет курс 1.0
	s.ratesCache[currency.RubleCurrency] = 1.0

	// Обновляем время истечения кэша
	s.cacheExpiry = time.Now().Add(s.cacheDuration)

	return nil
}

// isCacheValid проверяет, актуален ли кэш
func (s *ExchangeService) isCacheValid() bool {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()
	return time.Now().Before(s.cacheExpiry) && len(s.ratesCache) > 0
}

// calculateCommission рассчитывает комиссию
func (s *ExchangeService) calculateCommission(amount float64) float64 {
	commission := amount * (s.commissionRate / 100.0)
	// Округляем до 2 знаков после запятой
	return math.Round(commission*100) / 100
}
