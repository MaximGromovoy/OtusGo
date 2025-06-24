package mocks

import (
	"OtusGo/internal/interfaces"
	"OtusGo/internal/model/currency"
	"context"
	"fmt"
)

// MockCBRService реализует CBRServiceInterface
type MockCBRService struct {
	rates          map[string]float64
	supportedCurrs []string

	shouldFail bool
}

// Убеждаемся, что мок реализует интерфейс на этапе компиляции
var _ interfaces.CBRServiceInterface = (*MockCBRService)(nil)

// NewMockCBRService создает новый мок CBR сервиса
func NewMockCBRService() *MockCBRService {
	return &MockCBRService{
		rates: map[string]float64{
			currency.DollarCurrency: 90.0,
			currency.EuroCurrency:   95.0,
			currency.LiraCurrency:   2.8,
			currency.RubleCurrency:  1.0,
		},
		supportedCurrs: []string{
			currency.DollarCurrency,
			currency.EuroCurrency,
			currency.LiraCurrency,
			currency.RubleCurrency,
		},
	}
}

// GetCurrencyRates имитирует получение курсов валют
func (m *MockCBRService) GetCurrencyRates(ctx context.Context) ([]currency.CurrencyInterface, error) {

	if m.shouldFail {
		return nil, fmt.Errorf("mock error: failed to fetch currency rates")
	}

	var currencies []currency.CurrencyInterface
	for currType, rate := range m.rates {
		curr := currency.NewCurrency(currType, rate)
		if curr != nil {
			currencies = append(currencies, curr)
		}
	}

	return currencies, nil
}

// GetSupportedCurrencies возвращает список поддерживаемых валют
func (m *MockCBRService) GetSupportedCurrencies() []string {
	return m.supportedCurrs
}

// SetShouldFail позволяет задать поведение мока для тестов
func (m *MockCBRService) SetShouldFail(shouldFail bool) {
	m.shouldFail = shouldFail
}
