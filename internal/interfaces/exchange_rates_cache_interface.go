package interfaces

import "context"

// ExchangeRatesCacheInterface интерфейс для кэширования курсов валют
type ExchangeRatesCacheInterface interface {
	Set(ctx context.Context, fromCurrency, toCurrency string, rate float64) error
	Get(ctx context.Context, fromCurrency, toCurrency string) (float64, error)
	GetAll(ctx context.Context) (map[string]float64, error)
	Exists(ctx context.Context, fromCurrency, toCurrency string) (bool, error)
	Clear(ctx context.Context) error
	Close() error
}
