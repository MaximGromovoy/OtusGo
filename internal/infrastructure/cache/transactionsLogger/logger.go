package transactionsLogger

import (
	"OtusGo/internal/application/models/transaction"
	redisCache "OtusGo/internal/infrastructure/cache/redis"
	"context"
	"fmt"
	"time"
)

type TransactionsLogger struct {
	db *redisCache.RedisCache
}

func NewTransactionsLogger(db *redisCache.RedisCache) (*TransactionsLogger, error) {
	return &TransactionsLogger{
		db: db,
	}, nil
}

// LogTransaction логирует транзакцию обмена валют
func (l *TransactionsLogger) LogTransaction(ctx context.Context, ts *transaction.Transaction,
	duration time.Duration, errorMessage string) {
	transactionLog := newTransactionLog(ts, duration, errorMessage)
	logKey := fmt.Sprintf("operation:%d", transactionLog.TransactionID)
	l.db.Set(ctx, logKey, transactionLog, time.Minute)
}
