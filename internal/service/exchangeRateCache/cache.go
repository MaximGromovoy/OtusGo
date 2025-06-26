package exchangeRateCache

import (
	redisDatabase "OtusGo/internal/databases/redis"
	"context"
	"fmt"
	"time"
)

var ttl = time.Minute * 10

type ExchangeRatesCache struct {
	db *redisDatabase.RedisDatabase
}

func NewExchangeRatesCache(db *redisDatabase.RedisDatabase) (*ExchangeRatesCache, error) {
	return &ExchangeRatesCache{
		db: db,
	}, nil
}

// SetExchangeRate сохраняет курс обмена между двумя валютами
func (r *ExchangeRatesCache) Set(ctx context.Context, fromCurrency, toCurrency string, rate float64) error {
	return r.db.Set(ctx, r.buildRateKey(fromCurrency, toCurrency), rate, ttl)
}

// GetExchangeRate получает курс обмена между двумя валютами
func (r *ExchangeRatesCache) Get(ctx context.Context, fromCurrency, toCurrency string) (float64, error) {

	var rate float64
	err := r.db.Get(ctx, r.buildRateKey(fromCurrency, toCurrency), &rate)
	if err != nil {
		return 0, err
	}

	return rate, nil
}

// GetAll возвращает все закэшированные курсы валют
func (r *ExchangeRatesCache) GetAll(ctx context.Context) (map[string]float64, error) {
	keys, err := r.db.Keys(ctx, "rate:*")
	if err != nil {
		return nil, fmt.Errorf("failed to get rate keys: %w", err)
	}

	rates := make(map[string]float64)

	// Получаем значения для каждого ключа
	for _, key := range keys {
		var rate float64
		if err := r.db.Get(ctx, key, &rate); err != nil {
			continue
		}

		displayKey := key[5:]
		rates[displayKey] = rate
	}

	return rates, nil
}

// IsRateExpired проверяет, истек ли курс валют
func (r *ExchangeRatesCache) Exists(ctx context.Context, fromCurrency, toCurrency string) (bool, error) {
	return r.db.Exists(ctx, r.buildRateKey(fromCurrency, toCurrency))
}

func (r *ExchangeRatesCache) Clear(ctx context.Context) error {
	return r.db.DeleteByPattern(ctx, "*")
}

// Close закрывает соединение с Redis
func (r *ExchangeRatesCache) Close() error {
	return r.db.Close()
}

// buildRateKey создает ключ для курса валют
func (r *ExchangeRatesCache) buildRateKey(fromCurrency, toCurrency string) string {
	return fmt.Sprintf("rate:%s:%s", fromCurrency, toCurrency)
}
