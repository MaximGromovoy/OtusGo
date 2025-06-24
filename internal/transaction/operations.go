package transaction

import (
	"time"
)

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

// GetFormattedTimestamp возвращает отформатированное время
func (t *Transaction) GetFormattedTimestamp() string {
	return t.Timestamp.Format("2006-01-02 15:04:05")
}
