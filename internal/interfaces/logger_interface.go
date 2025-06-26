package interfaces

import (
	"OtusGo/internal/transaction"
	"context"
	"time"
)

// LoggerInterface интерфейс для логирования операций
type LoggerInterface interface {
	LogTransaction(ctx context.Context, tx *transaction.Transaction, duration time.Duration, errorMessage string)
}
