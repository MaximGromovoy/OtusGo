# Этап 1: Сборка (builder stage)
FROM golang:1.23.2-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости (кэшируется если go.mod не изменился)
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем приложение
# -o main - имя выходного файла
# ./cmd/server - путь к main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Этап 2: Финальный образ (production stage)
FROM alpine:latest

# Добавляем сертификаты для HTTPS запросов
RUN apk --no-cache add ca-certificates

# Создаем пользователя для безопасности
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Устанавливаем рабочую директорию
WORKDIR /home/appuser

# Копируем собранное приложение из builder stage
COPY --from=builder /app/main .

# Создаем директории для данных
RUN mkdir -p ./data && \
    chown -R appuser:appgroup /home/appuser

# Переключаемся на непривилегированного пользователя
USER appuser

# Открываем порт 8080
EXPOSE 8080

# Указываем команду для запуска приложения
CMD ["./main"]