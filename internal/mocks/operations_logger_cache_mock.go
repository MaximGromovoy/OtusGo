package mocks

import (
	"OtusGo/internal/interfaces"
	models "OtusGo/internal/model/operationLog"
	"context"
	"fmt"
	"sync"
)

// MockOperationsLoggerCache реализует OperationsLoggerCacheInterface
type MockOperationsLoggerCache struct {
	operations []models.OperationLog
	mutex      sync.RWMutex

	shouldFail  bool            // Флаг для имитации ошибок
	failMethods map[string]bool // Селективные ошибки для конкретных методов
}

var _ interfaces.OperationsLoggerCacheInterface = (*MockOperationsLoggerCache)(nil)

// NewMockOperationsLoggerCache создает новый мок кэша логирования операций
func NewMockOperationsLoggerCache() *MockOperationsLoggerCache {
	return &MockOperationsLoggerCache{
		operations:  make([]models.OperationLog, 0),
		failMethods: make(map[string]bool),
	}
}

// LogOperation имитирует логирование операции
func (m *MockOperationsLoggerCache) LogOperation(ctx context.Context, log *models.OperationLog) error {
	if m.shouldFail || m.failMethods["LogOperation"] {
		return fmt.Errorf("mock error: failed to log operation")
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Создаем копию для безопасности
	logCopy := models.OperationLog{
		TransactionID: log.TransactionID,
		UserID:        log.UserID,
		Operation:     log.Operation,
		Status:        log.Status,
		FromCurrency:  log.FromCurrency,
		ToCurrency:    log.ToCurrency,
		Amount:        log.Amount,
		Rate:          log.Rate,
		Commission:    log.Commission,
		Timestamp:     log.Timestamp,
		Duration:      log.Duration,
		ErrorMessage:  log.ErrorMessage,
	}

	m.operations = append(m.operations, logCopy)
	return nil
}

// Close имитирует закрытие соединения
func (m *MockOperationsLoggerCache) Close() error {
	if m.shouldFail || m.failMethods["Close"] {
		return fmt.Errorf("mock error: failed to close")
	}
	return nil
}

// Вспомогательные методы для настройки мока
func (m *MockOperationsLoggerCache) SetShouldFail(shouldFail bool) {
	m.shouldFail = shouldFail
}

func (m *MockOperationsLoggerCache) SetMethodFail(method string, shouldFail bool) {
	m.failMethods[method] = shouldFail
}
