package mocks

import (
	"OtusGo/internal/interfaces"
	"OtusGo/internal/transaction"
	"fmt"
	"sort"
	"sync"
)

// MockTransactionRepository реализует TransactionRepositoryInterface
type MockTransactionRepository struct {
	transactions map[int]*transaction.Transaction
	lastID       int
	mutex        sync.RWMutex

	shouldAddFail bool // Флаг для имитации ошибок
}

// Убеждаемся, что мок реализует интерфейс на этапе компиляции
var _ interfaces.TransactionRepositoryInterface = (*MockTransactionRepository)(nil)

// NewMockTransactionRepository создает новый мок репозитория
func NewMockTransactionRepository() *MockTransactionRepository {
	return &MockTransactionRepository{
		transactions: make(map[int]*transaction.Transaction),
		lastID:       0,
	}
}

// Add имитирует добавление транзакции
func (m *MockTransactionRepository) Add(tx *transaction.Transaction) error {

	if m.shouldAddFail {
		return fmt.Errorf("mock error: failed to add transaction")
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Увеличиваем lastID и присваиваем его транзакции
	m.lastID++
	tx.ID = m.lastID

	txCopy := &transaction.Transaction{
		ID:           tx.ID,
		UserID:       tx.UserID,
		Type:         tx.Type,
		Status:       tx.Status,
		FromCurrency: tx.FromCurrency,
		ToCurrency:   tx.ToCurrency,
		FromAmount:   tx.FromAmount,
		ToAmount:     tx.ToAmount,
		ExchangeRate: tx.ExchangeRate,
		Commission:   tx.Commission,
		Timestamp:    tx.Timestamp,
	}

	m.transactions[tx.ID] = txCopy

	return nil
}

// Get имитирует получение транзакции по ID
func (m *MockTransactionRepository) Get(id int) (*transaction.Transaction, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if tx, exists := m.transactions[id]; exists {
		txCopy := &transaction.Transaction{
			ID:           tx.ID,
			UserID:       tx.UserID,
			Type:         tx.Type,
			Status:       tx.Status,
			FromCurrency: tx.FromCurrency,
			ToCurrency:   tx.ToCurrency,
			FromAmount:   tx.FromAmount,
			ToAmount:     tx.ToAmount,
			ExchangeRate: tx.ExchangeRate,
			Commission:   tx.Commission,
			Timestamp:    tx.Timestamp,
		}
		return txCopy, nil
	}

	return nil, fmt.Errorf("transaction with ID %d not found", id)
}

// GetAll имитирует получение всех транзакций
func (m *MockTransactionRepository) GetAll() []*transaction.Transaction {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var result []*transaction.Transaction
	for _, tx := range m.transactions {
		txCopy := &transaction.Transaction{
			ID:           tx.ID,
			UserID:       tx.UserID,
			Type:         tx.Type,
			Status:       tx.Status,
			FromCurrency: tx.FromCurrency,
			ToCurrency:   tx.ToCurrency,
			FromAmount:   tx.FromAmount,
			ToAmount:     tx.ToAmount,
			ExchangeRate: tx.ExchangeRate,
			Commission:   tx.Commission,
			Timestamp:    tx.Timestamp,
		}
		result = append(result, txCopy)
	}

	// Сортируем по времени
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.After(result[j].Timestamp)
	})

	return result
}

func (m *MockTransactionRepository) SetShouldAddFail(shouldAddFail bool) {
	m.shouldAddFail = shouldAddFail
}
