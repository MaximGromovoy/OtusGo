package interfaces

import "context"

type ExchangeRatesProviderInterface interface {
	GetRate(ctx context.Context, fromCurrency, toCurrency string) (float64, error)
}
