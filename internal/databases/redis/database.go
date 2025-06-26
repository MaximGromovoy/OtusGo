package redisDatabase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisDatabase struct {
	client *redis.Client
}

func NewRedisDatabase(config *Configuration) (*RedisDatabase, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         config.address,
		Password:     config.password,
		DB:           config.dbNumber,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &RedisDatabase{
		client: rdb,
	}, nil
}

// Get универсальный метод чтения из кэша
func (r *RedisDatabase) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("cache miss for key: %s", key)
		}
		return fmt.Errorf("failed to get cache for key %s: %v", key, err)
	}

	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		return fmt.Errorf("failed to unmarshal to destination: %v", err)
	}

	return nil
}

// Set универсальный метод записи в кэш
func (r *RedisDatabase) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %v", err)
	}

	err = r.client.Set(ctx, key, data, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache for key %s: %v", key, err)
	}

	return nil
}

// GetString получает строковое значение
func (r *RedisDatabase) GetString(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// SetString сохраняет строковое значение
func (r *RedisDatabase) SetString(ctx context.Context, key, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

// Delete удаляет ключ из кэша
func (r *RedisDatabase) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// DeleteByPattern удаляет ключи по паттерну
func (r *RedisDatabase) DeleteByPattern(ctx context.Context, pattern string) error {
	// Получаем все ключи по паттерну
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get keys by pattern %s: %w", pattern, err)
	}

	// Если ключи найдены, удаляем их
	if len(keys) > 0 {
		err = r.client.Del(ctx, keys...).Err()
		if err != nil {
			return fmt.Errorf("failed to delete keys: %w", err)
		}
	}

	return nil
}

// Exists проверяет существование ключа
func (r *RedisDatabase) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	return count > 0, err
}

// Keys возвращает ключи по паттерну
func (r *RedisDatabase) Keys(ctx context.Context, pattern string) ([]string, error) {
	return r.client.Keys(ctx, pattern).Result()
}

// Close закрывает соединение с Redis
func (r *RedisDatabase) Close() error {
	return r.client.Close()
}

// Ping проверяет соединение с Redis
func (r *RedisDatabase) Ping(ctx context.Context) error {
	_, err := r.client.Ping(ctx).Result()
	return err
}
