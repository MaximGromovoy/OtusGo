package logger

import (
	redisDatabase "OtusGo/internal/databases/redis"
	"OtusGo/internal/interfaces"
	"OtusGo/internal/transaction"
	"context"
	"fmt"
	"time"
)

type Logger struct {
	db *redisDatabase.RedisDatabase
}

var _ interfaces.LoggerInterface = (*Logger)(nil)

func NewLogger(db *redisDatabase.RedisDatabase) (*Logger, error) {
	return &Logger{
		db: db,
	}, nil
}

// LogTransaction логирует транзакцию обмена валют
func (l *Logger) LogTransaction(ctx context.Context, tx *transaction.Transaction, duration time.Duration, errorMessage string) {
	logEntry := log{
		TransactionID: tx.ID,
		UserID:        tx.UserID,
		Operation:     string(tx.Type),
		Status:        string(tx.Status),
		FromCurrency:  tx.FromCurrency,
		ToCurrency:    tx.ToCurrency,
		Amount:        tx.ToAmount,
		Rate:          tx.ExchangeRate,
		Commission:    tx.Commission,
		Timestamp:     tx.Timestamp,
		Duration:      duration,
		ErrorMessage:  errorMessage,
	}

	logKey := fmt.Sprintf("operation:%d", tx.ID)
	l.db.Set(ctx, logKey, logEntry, time.Minute)
}
