package mocks

import (
	"OtusGo/internal/interfaces"
	"OtusGo/internal/transaction"
	"context"
	"time"
)

// MockLogger реализует LoggerInterface
type MockLogger struct{}

var _ interfaces.LoggerInterface = (*MockLogger)(nil)

// NewMockLogger создает новый мок кэша логирования операций
func NewMockLogger() *MockLogger {
	return &MockLogger{}
}

// LogTransaction имитирует логирование операции
func (m *MockLogger) LogTransaction(ctx context.Context,
	tx *transaction.Transaction, duration time.Duration, errorMessage string) {
}
