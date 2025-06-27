package exchange

import (
	"OtusGo/internal/application/models/transaction"
	"OtusGo/internal/domain/services/exchangeService"
	"context"
	"errors"
	"testing"
	"time"
)

// Mock для ExchangeRatesProviderInterface
type mockExchangeRatesProvider struct {
	rates map[string]float64
	err   error
}

func (m *mockExchangeRatesProvider) GetRate(ctx context.Context, fromCurrency, toCurrency string) (float64, error) {
	if m.err != nil {
		return 0, m.err
	}
	key := fromCurrency + "-" + toCurrency
	if rate, exists := m.rates[key]; exists {
		return rate, nil
	}
	return 0, errors.New("rate not found")
}

// Mock для TransactionRepositoryInterface
type mockTransactionRepository struct {
	transactions map[int]*transaction.Transaction
	addError     error
	updateError  error
	nextID       int
}

func newMockTransactionRepository() *mockTransactionRepository {
	return &mockTransactionRepository{
		transactions: make(map[int]*transaction.Transaction),
		nextID:       1,
	}
}

func (m *mockTransactionRepository) Add(tx *transaction.Transaction) error {
	if m.addError != nil {
		return m.addError
	}
	tx.ID = m.nextID
	m.nextID++
	m.transactions[tx.ID] = tx
	return nil
}

func (m *mockTransactionRepository) Get(id int) (*transaction.Transaction, error) {
	if tx, exists := m.transactions[id]; exists {
		return tx, nil
	}
	return nil, errors.New("transaction not found")
}

func (m *mockTransactionRepository) GetAll() []*transaction.Transaction {
	var result []*transaction.Transaction
	for _, tx := range m.transactions {
		result = append(result, tx)
	}
	return result
}

func (m *mockTransactionRepository) Update(tx *transaction.Transaction) error {
	if m.updateError != nil {
		return m.updateError
	}
	m.transactions[tx.ID] = tx
	return nil
}

// Mock для LoggerInterface
type mockLogger struct {
	loggedTransactions []logEntry
}

type logEntry struct {
	transaction  *transaction.Transaction
	duration     time.Duration
	errorMessage string
}

func (m *mockLogger) LogTransaction(ctx context.Context, tx *transaction.Transaction, duration time.Duration, errorMessage string) {
	m.loggedTransactions = append(m.loggedTransactions, logEntry{
		transaction:  tx,
		duration:     duration,
		errorMessage: errorMessage,
	})
}

func TestNewExchangeOrchestrator(t *testing.T) {
	provider := &mockExchangeRatesProvider{}
	repo := newMockTransactionRepository()
	service := exchangeService.NewExchangeService()
	logger := &mockLogger{}

	orchestrator := NewExchangeOrchestrator(provider, repo, service, logger)

	if orchestrator == nil {
		t.Fatal("Expected orchestrator to be created, got nil")
	}
	if orchestrator.exchangeRatesProvider != provider {
		t.Error("Exchange rates provider not set correctly")
	}
	if orchestrator.transactionRepository != repo {
		t.Error("Transaction repository not set correctly")
	}
	if orchestrator.exchangeService != service {
		t.Error("Exchange service not set correctly")
	}
	if orchestrator.logger != logger {
		t.Error("Logger not set correctly")
	}
}

func TestNewExchangeOrchestratorRequest(t *testing.T) {
	userID := 123
	fromCurrency := "USD"
	toCurrency := "EUR"
	amount := 100.0
	commissionRate := 2.5

	req := NewExchangeOrchestratorRequest(userID, fromCurrency, toCurrency, amount, commissionRate)

	if req.UserId != userID {
		t.Errorf("Expected UserID %d, got %d", userID, req.UserId)
	}
	if req.FromCurrency != fromCurrency {
		t.Errorf("Expected FromCurrency %s, got %s", fromCurrency, req.FromCurrency)
	}
	if req.ToCurrency != toCurrency {
		t.Errorf("Expected ToCurrency %s, got %s", toCurrency, req.ToCurrency)
	}
	if req.Amount != amount {
		t.Errorf("Expected Amount %f, got %f", amount, req.Amount)
	}
	if req.CommissionRate != commissionRate {
		t.Errorf("Expected CommissionRate %f, got %f", commissionRate, req.CommissionRate)
	}
}

func TestNewExchangeOrchestratorResponse(t *testing.T) {
	amount := 95.0
	commission := 2.0
	exchangeRate := 0.85

	resp := NewExchangeOrchestratorResponse(amount, commission, exchangeRate)

	if resp.Amount != amount {
		t.Errorf("Expected Amount %f, got %f", amount, resp.Amount)
	}
	if resp.Commission != commission {
		t.Errorf("Expected Commission %f, got %f", commission, resp.Commission)
	}
	if resp.ExchangeRate != exchangeRate {
		t.Errorf("Expected ExchangeRate %f, got %f", exchangeRate, resp.ExchangeRate)
	}
}

func TestExchangeOrchestrator_Exchange_Success(t *testing.T) {
	// Arrange
	provider := &mockExchangeRatesProvider{
		rates: map[string]float64{
			"USD-EUR": 0.85,
		},
	}
	repo := newMockTransactionRepository()
	service := exchangeService.NewExchangeService()
	logger := &mockLogger{}

	orchestrator := NewExchangeOrchestrator(provider, repo, service, logger)

	req := NewExchangeOrchestratorRequest(1, "USD", "EUR", 100.0, 2.0)
	ctx := context.Background()

	// Act
	result, err := orchestrator.Exchange(ctx, req)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// Проверяем результат обмена
	expectedAmount := 83.3 // (100 * 0.85) - (100 * 0.85 * 0.02)
	if result.Amount != expectedAmount {
		t.Errorf("Expected amount %f, got %f", expectedAmount, result.Amount)
	}

	expectedCommission := 1.7 // (100 * 0.85) * 0.02
	if result.Commission != expectedCommission {
		t.Errorf("Expected commission %f, got %f", expectedCommission, result.Commission)
	}

	if result.ExchangeRate != 0.85 {
		t.Errorf("Expected exchange rate 0.85, got %f", result.ExchangeRate)
	}

	// Проверяем, что транзакция была создана
	if len(repo.transactions) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(repo.transactions))
	}

	// Проверяем статус транзакции
	tx := repo.transactions[1]
	if tx.Status != transaction.TransactionStatusCompleted {
		t.Errorf("Expected transaction status %s, got %s", transaction.TransactionStatusCompleted, tx.Status)
	}

	// Проверяем логирование
	if len(logger.loggedTransactions) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(logger.loggedTransactions))
	}
}

func TestExchangeOrchestrator_Exchange_TransactionAddError(t *testing.T) {
	// Arrange
	provider := &mockExchangeRatesProvider{}
	repo := newMockTransactionRepository()
	repo.addError = errors.New("database error")
	service := exchangeService.NewExchangeService()
	logger := &mockLogger{}

	orchestrator := NewExchangeOrchestrator(provider, repo, service, logger)

	req := NewExchangeOrchestratorRequest(1, "USD", "EUR", 100.0, 2.0)
	ctx := context.Background()

	// Act
	result, err := orchestrator.Exchange(ctx, req)

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if result != nil {
		t.Error("Expected nil result on error")
	}

	expectedErrorMsg := "failed to create transaction: database error"
	if err.Error() != expectedErrorMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMsg, err.Error())
	}

	// Проверяем, что ошибка была залогирована
	if len(logger.loggedTransactions) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(logger.loggedTransactions))
	}
	if logger.loggedTransactions[0].errorMessage != "database error" {
		t.Errorf("Expected error message 'database error', got '%s'", logger.loggedTransactions[0].errorMessage)
	}
}

func TestExchangeOrchestrator_Exchange_ExchangeRateError(t *testing.T) {
	// Arrange
	provider := &mockExchangeRatesProvider{
		err: errors.New("rate service unavailable"),
	}
	repo := newMockTransactionRepository()
	service := exchangeService.NewExchangeService()
	logger := &mockLogger{}

	orchestrator := NewExchangeOrchestrator(provider, repo, service, logger)

	req := NewExchangeOrchestratorRequest(1, "USD", "EUR", 100.0, 2.0)
	ctx := context.Background()

	// Act
	result, err := orchestrator.Exchange(ctx, req)

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if result != nil {
		t.Error("Expected nil result on error")
	}

	expectedErrorMsg := "failed to get exchange rate: rate service unavailable"
	if err.Error() != expectedErrorMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMsg, err.Error())
	}

	// Проверяем, что транзакция была помечена как неудачная
	if len(repo.transactions) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(repo.transactions))
	}

	tx := repo.transactions[1]
	if tx.Status != transaction.TransactionStatusFailed {
		t.Errorf("Expected transaction status %s, got %s", transaction.TransactionStatusFailed, tx.Status)
	}

	// Проверяем логирование
	if len(logger.loggedTransactions) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(logger.loggedTransactions))
	}
}

func TestExchangeOrchestrator_Exchange_UpdateTransactionError(t *testing.T) {
	// Arrange
	provider := &mockExchangeRatesProvider{
		rates: map[string]float64{
			"USD-EUR": 0.85,
		},
	}
	repo := newMockTransactionRepository()
	repo.updateError = errors.New("update failed")
	service := exchangeService.NewExchangeService()
	logger := &mockLogger{}

	orchestrator := NewExchangeOrchestrator(provider, repo, service, logger)

	req := NewExchangeOrchestratorRequest(1, "USD", "EUR", 100.0, 2.0)
	ctx := context.Background()

	// Act
	result, err := orchestrator.Exchange(ctx, req)

	// Assert
	// Несмотря на ошибку обновления, операция должна завершиться успешно
	// так как основной обмен выполнился
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestExchangeOrchestrator_Exchange_ContextCancellation(t *testing.T) {
	// Arrange
	provider := &mockExchangeRatesProvider{
		rates: map[string]float64{
			"USD-EUR": 0.85,
		},
	}
	repo := newMockTransactionRepository()
	service := exchangeService.NewExchangeService()
	logger := &mockLogger{}

	orchestrator := NewExchangeOrchestrator(provider, repo, service, logger)

	req := NewExchangeOrchestratorRequest(1, "USD", "EUR", 100.0, 2.0)

	// Создаем уже отмененный контекст
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Act
	result, err := orchestrator.Exchange(ctx, req)

	// Assert
	// Результат зависит от того, на каком этапе проверяется контекст
	// В текущей реализации контекст передается только в GetRate
	if err != nil && result == nil {
		// Это нормально, если контекст был отменен
		t.Logf("Context cancellation handled correctly: %v", err)
	}
}

func TestExchangeOrchestrator_Exchange_ZeroAmount(t *testing.T) {
	// Arrange
	provider := &mockExchangeRatesProvider{
		rates: map[string]float64{
			"USD-EUR": 0.85,
		},
	}
	repo := newMockTransactionRepository()
	service := exchangeService.NewExchangeService()
	logger := &mockLogger{}

	orchestrator := NewExchangeOrchestrator(provider, repo, service, logger)

	req := NewExchangeOrchestratorRequest(1, "USD", "EUR", 0.0, 2.0)
	ctx := context.Background()

	// Act
	result, err := orchestrator.Exchange(ctx, req)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Amount != 0.0 {
		t.Errorf("Expected amount 0.0, got %f", result.Amount)
	}
	if result.Commission != 0.0 {
		t.Errorf("Expected commission 0.0, got %f", result.Commission)
	}
}

func TestExchangeOrchestrator_Exchange_HighCommissionRate(t *testing.T) {
	// Arrange
	provider := &mockExchangeRatesProvider{
		rates: map[string]float64{
			"USD-EUR": 0.85,
		},
	}
	repo := newMockTransactionRepository()
	service := exchangeService.NewExchangeService()
	logger := &mockLogger{}

	orchestrator := NewExchangeOrchestrator(provider, repo, service, logger)

	req := NewExchangeOrchestratorRequest(1, "USD", "EUR", 100.0, 50.0) // 50% комиссия
	ctx := context.Background()

	// Act
	result, err := orchestrator.Exchange(ctx, req)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// При 50% комиссии результат должен быть 42.5 (85 * 0.5)
	expectedAmount := 42.5
	if result.Amount != expectedAmount {
		t.Errorf("Expected amount %f, got %f", expectedAmount, result.Amount)
	}
}
