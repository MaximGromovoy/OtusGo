package transactionsLogger

import (
	"OtusGo/internal/application/models/transaction"
	"time"
)

type transactionLog struct {
	TransactionID int           `json:"transaction_id"`
	UserID        int           `json:"user_id"`
	Operation     string        `json:"operation"`
	Status        string        `json:"status"`
	FromCurrency  string        `json:"from_currency"`
	ToCurrency    string        `json:"to_currency"`
	Amount        float64       `json:"amount"`
	Rate          float64       `json:"rate"`
	Commission    float64       `json:"commission"`
	Timestamp     time.Time     `json:"timestamp"`
	Duration      time.Duration `json:"duration"`
	ErrorMessage  string        `json:"error_message"`
}

func newTransactionLog(ts *transaction.Transaction,
	duration time.Duration, errorMessage string) transactionLog {
	return transactionLog{
		TransactionID: ts.ID,
		UserID:        ts.UserID,
		Operation:     string(ts.Type),
		Status:        string(ts.Status),
		FromCurrency:  ts.FromCurrency,
		ToCurrency:    ts.ToCurrency,
		Amount:        ts.ToAmount,
		Rate:          ts.ExchangeRate,
		Commission:    ts.Commission,
		Timestamp:     ts.Timestamp,
		Duration:      duration,
		ErrorMessage:  errorMessage,
	}
}
