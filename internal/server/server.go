// Package server предоставляет HTTP API для работы с валютами
// @title API Валют
// @version 1.0
// @description API для управления различными типами валют
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey JWT
// @in header
// @name Authorization
package server

import (
	"OtusGo/internal/repository/currencyRepository"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

// CurrencyServer представляет сервер для обработки запросов по валютам
type CurrencyServer struct {
	repo      *currencyRepository.CurrencyRepository
	jwtSecret string
}

// NewCurrencyServer создает и возвращает новый экземпляр CurrencyServer
func NewCurrencyServer(repo *currencyRepository.CurrencyRepository) *CurrencyServer {
	return &CurrencyServer{
		repo:      repo,
		jwtSecret: "secret-key",
	}
}

// Start инициализирует и запускает HTTP сервер
func (s *CurrencyServer) Start(port string) error {
	mux := http.NewServeMux()

	// Регистрация обработчиков
	mux.HandleFunc("/login", s.handleLogin)
	mux.HandleFunc("/currency", s.handleCurrency)
	mux.HandleFunc("/currencies", s.handleCurrencies)
	mux.HandleFunc("/currency/", s.handleCurrencyWithTypeAndID)

	// Добавление Swagger UI
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	return http.ListenAndServe(":"+port, mux)
}
