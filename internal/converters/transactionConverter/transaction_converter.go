package transactionConverter

import (
	"OtusGo/internal/transaction"
	"fmt"
	"strconv"
	"time"
)

var headers = []string{
	"id", "user_id", "type", "status", "from_currency",
	"to_currency", "from_amount", "to_amount",
	"exchange_rate", "commission", "timestamp",
}

func ToRecord(transaction *transaction.Transaction) ([]string, error) {
	var record = []string{
		strconv.Itoa(transaction.ID),
		strconv.Itoa(transaction.UserID),
		string(transaction.Type),
		string(transaction.Status),
		transaction.FromCurrency,
		transaction.ToCurrency,
		fmt.Sprintf("%.6f", transaction.FromAmount),
		fmt.Sprintf("%.6f", transaction.ToAmount),
		fmt.Sprintf("%.6f", transaction.ExchangeRate),
		fmt.Sprintf("%.6f", transaction.Commission),
		transaction.Timestamp.Format(time.RFC3339),
	}

	if len(record) != len(headers) {
		return nil, fmt.Errorf("invalid record length: expected %d, got %d", len(headers), len(record))
	}

	return record, nil
}

func ToModel(record []string) (*transaction.Transaction, error) {
	if len(record) != len(headers) {
		return nil, fmt.Errorf("invalid record length: expected %d, got %d", len(headers), len(record))
	}

	id, err := strconv.Atoi(record[0])
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %v", err)
	}

	userId, err := strconv.Atoi(record[1])
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %v", err)
	}

	var transactionType = transaction.TransactionType(record[2])
	var transactionStatus = transaction.TransactionStatus(record[3])

	var fromCurrency = record[4]
	var toCurrency = record[5]

	fromAmount, err := strconv.ParseFloat(record[6], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid from_amount: %v", err)
	}

	toAmount, err := strconv.ParseFloat(record[7], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid to_amount: %v", err)
	}

	exchangeRate, err := strconv.ParseFloat(record[8], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid exchange_rate: %v", err)
	}

	commission, err := strconv.ParseFloat(record[9], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid commission: %v", err)
	}

	timestamp, err := time.Parse(time.RFC3339, record[10])
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp: %v", err)
	}

	return &transaction.Transaction{
		ID:           id,
		UserID:       userId,
		Type:         transactionType,
		Status:       transactionStatus,
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		FromAmount:   fromAmount,
		ToAmount:     toAmount,
		ExchangeRate: exchangeRate,
		Commission:   commission,
		Timestamp:    timestamp,
	}, nil
}

func GetHeaders() []string {
	return headers
}
