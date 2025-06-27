package transaction

import (
	"time"
)

// TransactionType определяет тип транзакции
type TransactionType string

const (
	TransactionTypeExchange TransactionType = "exchange" // Обмен валют
)

// TransactionStatus определяет статус транзакции
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"   // В процессе
	TransactionStatusCompleted TransactionStatus = "completed" // Завершена
	TransactionStatusFailed    TransactionStatus = "failed"    // Ошибка
)

// Transaction представляет запись о транзакции
type Transaction struct {
	ID           int               // Уникальный идентификатор транзакции
	UserID       int               // ID пользователя
	Type         TransactionType   // Тип операции
	Status       TransactionStatus // Статус
	FromCurrency string            // Исходная валюта
	ToCurrency   string            // Целевая валюта
	FromAmount   float64           // Исходная сумма
	ToAmount     float64           // Итоговая сумма
	ExchangeRate float64           // Курс обмена
	Commission   float64           // Комиссия
	Timestamp    time.Time         // Время операции
}

// NewExchangeTransaction создает транзакцию обмена валют
func NewExchangeTransaction(userID int, fromCurrency string, toCurrency string,
	fromAmount float64) *Transaction {
	return &Transaction{
		UserID:       userID,
		Type:         TransactionTypeExchange,
		Status:       TransactionStatusPending,
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		FromAmount:   fromAmount,
		Timestamp:    time.Now(),
	}
}

// Complete помечает транзакцию как завершенную
func (t *Transaction) Complete() {
	t.Status = TransactionStatusCompleted
}

// Fail помечает транзакцию как неудачную
func (t *Transaction) Fail() {
	t.Status = TransactionStatusFailed
}

// GetFormattedTimestamp возвращает отформатированное время
func (t *Transaction) GetFormattedTimestamp() string {
	return t.Timestamp.Format("2006-01-02 15:04:05")
}
