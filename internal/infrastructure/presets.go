package presets

import (
	redisCache "OtusGo/internal/infrastructure/cache/redis"
	postgresRepository "OtusGo/internal/infrastructure/repository/postgres"
	"os"
	"strconv"
)

func GetPostgresConfigurationPreset() *postgresRepository.Configuration {
	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5432")
	user := getEnv("POSTGRES_USER", "postgres")
	password := getEnv("POSTGRES_PASSWORD", "password")
	dbName := getEnv("POSTGRES_DB", "exchange_service_db")
	sslMode := getEnv("POSTGRES_SSLMODE", "disable")

	postgresUri := "postgres://" + user + ":" + password + "@" + host + ":" + port

	return postgresRepository.NewConfigurationFromParts(postgresUri, dbName, "sslmode="+sslMode)
}

func GetRedisConfigurationPreset() *redisCache.Configuration {
	host := getEnv("REDIS_HOST", "localhost")
	port := getEnv("REDIS_PORT", "6379")
	password := getEnv("REDIS_PASSWORD", "")
	dbNum, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	address := host + ":" + port

	return redisCache.NewConfiguration(address, password, dbNum)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
