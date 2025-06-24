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
