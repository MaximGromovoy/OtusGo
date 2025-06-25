package interfaces

import (
	models "OtusGo/internal/model/operationLog"
	"context"
)

// OperationsLoggerCacheInterface интерфейс для логирования операций
type OperationsLoggerCacheInterface interface {
	LogOperation(ctx context.Context, log *models.OperationLog) error
	Close() error
}
