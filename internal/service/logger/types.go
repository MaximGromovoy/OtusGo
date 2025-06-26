package logger

import "time"

type log struct {
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
