package mocks

import (
	"OtusGo/internal/interfaces"
	"context"
	"fmt"
	"sync"
)

// MockExchangeRatesCache реализует ExchangeRatesCacheInterface
type MockExchangeRatesCache struct {
	rates       map[string]float64 // ключ формата "FROM:TO"
	mutex       sync.RWMutex
	shouldFail  bool            // Флаг для имитации ошибок
	failMethods map[string]bool // Селективные ошибки для конкретных методов
}

var _ interfaces.ExchangeRatesCacheInterface = (*MockExchangeRatesCache)(nil)

// NewMockExchangeRatesCache создает новый мок кэша курсов валют
func NewMockExchangeRatesCache() *MockExchangeRatesCache {
	return &MockExchangeRatesCache{
		rates:       make(map[string]float64),
		failMethods: make(map[string]bool),
	}
}

// Set имитирует сохранение курса валют
func (m *MockExchangeRatesCache) Set(ctx context.Context, fromCurrency, toCurrency string, rate float64) error {
	if m.shouldFail || m.failMethods["Set"] {
		return fmt.Errorf("mock error: failed to set rate")
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	key := fmt.Sprintf("%s:%s", fromCurrency, toCurrency)
	m.rates[key] = rate
	return nil
}

// Get имитирует получение курса валют
func (m *MockExchangeRatesCache) Get(ctx context.Context, fromCurrency, toCurrency string) (float64, error) {
	if m.shouldFail || m.failMethods["Get"] {
		return 0, fmt.Errorf("mock error: failed to get rate")
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	key := fmt.Sprintf("%s:%s", fromCurrency, toCurrency)
	rate, exists := m.rates[key]
	if !exists {
		return 0, fmt.Errorf("cache miss for key: %s", key)
	}

	return rate, nil
}

// GetAll имитирует получение всех курсов валют
func (m *MockExchangeRatesCache) GetAll(ctx context.Context) (map[string]float64, error) {
	if m.shouldFail || m.failMethods["GetAll"] {
		return nil, fmt.Errorf("mock error: failed to get all rates")
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make(map[string]float64)
	for key, rate := range m.rates {
		result[key] = rate
	}

	return result, nil
}

// Exists имитирует проверку существования курса валют
func (m *MockExchangeRatesCache) Exists(ctx context.Context, fromCurrency, toCurrency string) (bool, error) {
	if m.shouldFail || m.failMethods["Exists"] {
		return false, fmt.Errorf("mock error: failed to check existence")
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	key := fmt.Sprintf("%s:%s", fromCurrency, toCurrency)
	_, exists := m.rates[key]
	return exists, nil
}

// Clear имитирует очистку кэша
func (m *MockExchangeRatesCache) Clear(ctx context.Context) error {
	if m.shouldFail || m.failMethods["Clear"] {
		return fmt.Errorf("mock error: failed to clear cache")
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.rates = make(map[string]float64)
	return nil
}

// Close имитирует закрытие соединения
func (m *MockExchangeRatesCache) Close() error {
	if m.shouldFail || m.failMethods["Close"] {
		return fmt.Errorf("mock error: failed to close")
	}
	return nil
}

// Вспомогательные методы для настройки мока
func (m *MockExchangeRatesCache) SetShouldFail(shouldFail bool) {
	m.shouldFail = shouldFail
}

func (m *MockExchangeRatesCache) SetMethodFail(method string, shouldFail bool) {
	m.failMethods[method] = shouldFail
}

func (m *MockExchangeRatesCache) GetRatesCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.rates)
}

func (m *MockExchangeRatesCache) HasRate(fromCurrency, toCurrency string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	key := fmt.Sprintf("%s:%s", fromCurrency, toCurrency)
	_, exists := m.rates[key]
	return exists
}
