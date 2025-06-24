package interfaces

import (
	"OtusGo/internal/model/currency"
	"context"
)

// CBRServiceInterface определяет контракт для сервиса курсов валют
type CBRServiceInterface interface {
	GetCurrencyRates(ctx context.Context) ([]currency.CurrencyInterface, error)
	GetSupportedCurrencies() []string
}

// CurrencyServiceInterface - общий интерфейс для сервисов курсов валют
type CurrencyServiceInterface interface {
	CBRServiceInterface
}
