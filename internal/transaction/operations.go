package transaction

import (
	"time"
)

// NewTransaction создает новую транзакцию
func NewTransaction(userID int, txType TransactionType, toCurrency string, toAmount float64, description string) *Transaction {
	return &Transaction{
		UserID:     userID,
		Type:       txType,
		Status:     TransactionStatusPending,
		ToCurrency: toCurrency,
		ToAmount:   toAmount,
		Timestamp:  time.Now(),
	}
}

// NewDepositTransaction создает транзакцию пополнения
func NewDepositTransaction(userID int, fromCurrency string, toCurrency string,
	fromAmount, toAmount, exchangeRate, commission float64) *Transaction {
	return &Transaction{
		UserID:       userID,
		Type:         TransactionTypeDeposit,
		Status:       TransactionStatusPending,
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		FromAmount:   fromAmount,
		ToAmount:     toAmount,
		ExchangeRate: exchangeRate,
		Commission:   commission,
		Timestamp:    time.Now(),
	}
}

// NewWithdrawTransaction создает транзакцию снятия
func NewWithdrawTransaction(userID int, fromCurrency string, toCurrency string,
	fromAmount, toAmount, exchangeRate, commission float64) *Transaction {
	return &Transaction{
		UserID:       userID,
		Type:         TransactionTypeWithdraw,
		Status:       TransactionStatusPending,
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		FromAmount:   fromAmount,
		ToAmount:     toAmount,
		ExchangeRate: exchangeRate,
		Commission:   commission,
		Timestamp:    time.Now(),
	}
}

// NewExchangeTransaction создает транзакцию обмена валют
func NewExchangeTransaction(userID int, fromCurrency string, toCurrency string,
	fromAmount, toAmount, exchangeRate, commission float64) *Transaction {
	return &Transaction{
		UserID:       userID,
		Type:         TransactionTypeExchange,
		Status:       TransactionStatusPending,
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		FromAmount:   fromAmount,
		ToAmount:     toAmount,
		ExchangeRate: exchangeRate,
		Commission:   commission,
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

// IsCompleted проверяет, завершена ли транзакция
func (t *Transaction) IsCompleted() bool {
	return t.Status == TransactionStatusCompleted
}

// GetFormattedTimestamp возвращает отформатированное время
func (t *Transaction) GetFormattedTimestamp() string {
	return t.Timestamp.Format("2006-01-02 15:04:05")
}
