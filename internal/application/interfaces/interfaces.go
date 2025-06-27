package interfaces

import (
	"OtusGo/internal/application/models/transaction"
	"OtusGo/internal/domain/models/currency"
	"context"
	"time"
)

// LoggerInterface интерфейс для логирования операций
type LoggerInterface interface {
	LogTransaction(ctx context.Context, tx *transaction.Transaction, duration time.Duration, errorMessage string)
}

// TransactionRepositoryInterface определяет контракт для репозитория транзакций
type TransactionRepositoryInterface interface {
	Add(tx *transaction.Transaction) error
	Get(id int) (*transaction.Transaction, error)
	GetAll() []*transaction.Transaction
	Update(tx *transaction.Transaction) error
}

// ExchangeRatesProviderInterface определяет контракт получения курсов обмена валют
type ExchangeRatesProviderInterface interface {
	GetRate(ctx context.Context, fromCurrency, toCurrency string) (float64, error)
}

// CBRServiceInterface определяет контракт для сервиса курсов валют
type CBRServiceInterface interface {
	GetCurrencyRates(ctx context.Context) ([]*currency.Currency, error)
	GetSupportedCurrencies(ctx context.Context) ([]string, error)
}

// ExchangeRatesCacheInterface интерфейс для кэширования курсов валют
type ExchangeRatesCacheInterface interface {
	Set(ctx context.Context, fromCurrency, toCurrency string, rate float64) error
	Get(ctx context.Context, fromCurrency, toCurrency string) (float64, error)
	GetAll(ctx context.Context) (map[string]float64, error)
	Exists(ctx context.Context, fromCurrency, toCurrency string) (bool, error)
	Clear(ctx context.Context) error
	Close() error
}
