package operationsLoggerCache

import (
	models "OtusGo/internal/model/operationLog"
	"OtusGo/internal/service/redisCache"
	"context"
	"fmt"
	"time"
)

type OperationsLoggerCache struct {
	cache *redisCache.RedisCache
}

func NewOperationsLoggerCache(config redisCache.CacheConfig) (*OperationsLoggerCache, error) {
	cache, err := redisCache.NewRedisCache(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create operations logger cache: %v", err)
	}

	return &OperationsLoggerCache{
		cache: cache,
	}, nil
}

// LogOperation логирует операцию обмена валют
func (l *OperationsLoggerCache) LogOperation(ctx context.Context, log *models.OperationLog) error {
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}

	logKey := fmt.Sprintf("operation:%d", log.TransactionID)
	return l.cache.Set(ctx, logKey, log, time.Minute)
}

// Close закрывает соединение
func (l *OperationsLoggerCache) Close() error {
	return l.cache.Close()
}
