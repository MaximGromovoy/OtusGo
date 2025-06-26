package exchangeratesprovider

import (
	"OtusGo/internal/interfaces"
	"context"
	"fmt"
)

type ExchangeRatesProvider struct {
	cbrService         interfaces.CBRServiceInterface
	exchangeRatesCache interfaces.ExchangeRatesCacheInterface
}

func NewExchangeRatesProvider(cbrService interfaces.CBRServiceInterface,
	exchangeRatesCache interfaces.ExchangeRatesCacheInterface) *ExchangeRatesProvider {
	return &ExchangeRatesProvider{
		cbrService:         cbrService,
		exchangeRatesCache: exchangeRatesCache,
	}
}

// getExchangeRate получает курс обмена между двумя валютами
func (p *ExchangeRatesProvider) GetRate(ctx context.Context, fromCurrency, toCurrency string) (float64, error) {
	exist, err := p.exchangeRatesCache.Exists(ctx, fromCurrency, toCurrency)
	if err != nil {
		return 0, fmt.Errorf("failed to check if rate exists: %w", err)
	}

	// Если курс уже есть в кэше, возвращаем его
	if exist {
		rate, err := p.exchangeRatesCache.Get(ctx, fromCurrency, toCurrency)
		if err != nil {
			return 0, fmt.Errorf("failed to get cached rate: %w", err)
		}
		return rate, nil
	}

	// Если курс не найден в кэше, получаем его из CBR
	currencies, err := p.cbrService.GetCurrencyRates(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch currency rates: %w", err)
	}

	var fromRate, toRate float64
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
	}

	if fromRate == 0 || toRate == 0 {
		return 0, fmt.Errorf("currency rates not found for %s or %s", fromCurrency, toCurrency)
	}

	// Нормализуем курсы относительно рубля
	fromRateToRub := fromRate
	toRateToRub := toRate

	if fromCurrency == "RUB" {
		fromRateToRub = 1.0
	}
	if toCurrency == "RUB" {
		toRateToRub = 1.0
	}

	// Рассчитываем прямой курс обмена
	exchangeRate := fromRateToRub / toRateToRub

	if err := p.exchangeRatesCache.Set(ctx, fromCurrency, toCurrency, exchangeRate); err != nil {
		return 0, fmt.Errorf("failed to cache exchange rate: %w", err)
	}

	return exchangeRate, nil
}
