# OtusGo - Сервис обмена валют

Микросервис для обмена валют с получением актуальных курсов от ЦБ РФ, реализованный на Go с использованием принципов Clean Architecture.

## Возможности

- **Обмен валют** с актуальными курсами ЦБ РФ
- **HTTP API** для интеграции с другими сервисами
- **Консольное приложение** для интерактивного использования
- **Кэширование курсов** в Redis для повышения производительности
- **Логирование транзакций** с метриками времени выполнения
- **Retry механизм** для устойчивости к сбоям внешних API
- **Docker контейнеризация** для простого развертывания

## Архитектура

```
├── cmd/                    # Точки входа приложения
│   ├── console/           # Консольное приложение
│   └── server/            # HTTP сервер
├── internal/
│   ├── application/       # Слой приложения
│   │   ├── interfaces/    # Интерфейсы
│   │   ├── models/        # Модели приложения
│   │   └── services/      # Сервисы приложения
│   ├── domain/            # Доменный слой
│   │   ├── models/        # Доменные модели
│   │   └── services/      # Доменные сервисы
│   └── infrastructure/    # Слой инфраструктуры
│       ├── cache/         # Кэширование (Redis)
│       ├── external/      # Внешние API (ЦБ РФ)
│       ├── httpServer/    # HTTP сервер и хендлеры
│       └── repository/    # Репозитории (PostgreSQL)
```

## 📦 Установка и запуск

### Запуск через Docker Compose

```bash
# Запуск инфраструктуры (PostgreSQL, Redis, админки)
docker-compose up postgres redis redis-commander pgadmin

# Запуск HTTP сервера
docker-compose --profile server up --build

# Запуск консольного приложения
docker-compose --profile console up --build
```

## API Endpoints

### HTTP API

```bash
# Получить курс валют
GET /api/rates?from=USD&to=EUR

# Обменять валюту
POST /api/exchange
Content-Type: application/json

{
    "user_id": 1,
    "from_currency": "USD",
    "to_currency": "EUR",
    "amount": 100.0,
    "commission_rate": 2.0
}
```

### Консольные команды

```bash
# Обмен валют
exchange

# Получить курс
rate

# Выход
exit
```

## Тестирование

```bash
# Запуск тестов
go test ./internal/application/services/exchange -v
```

## Мониторинг

### Веб-интерфейсы

- **pgAdmin**: http://localhost:5050
  - Email: admin@admin.com
  - Password: admin

- **Redis Commander**: http://localhost:8082

### Логирование

Приложение использует структурированное логирование с slog:
- Все HTTP запросы логируются
- Транзакции обмена валют логируются с метриками
- Ошибки внешних API логируются с retry попытками

### Docker Compose профили

- `server` - HTTP сервер
- `console` - консольное приложение

## 🔍 Примеры использования

### Получение курса через HTTP API

```bash
curl -X GET "http://localhost:8080/api/rates?from=USD&to=EUR"
```

### Обмен валюты через HTTP API

```bash
curl -X POST "http://localhost:8080/api/exchange" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "from_currency": "USD",
    "to_currency": "EUR",
    "amount": 100.0,
    "commission_rate": 2.0
  }'
```

### Использование консольного приложения

```bash
# Запуск консольного приложения
docker-compose --profile console up

# Интерактивное использование
Введите команду: exchange
Введите исходную валюту (например, USD): USD
Введите целевую валюту (например, EUR): EUR
Введите сумму для обмена: 100
Введите процент комиссии (например, 2.5): 2.0
```