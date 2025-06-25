package server

import (
	"OtusGo/internal/httpServer/handlers"
	"context"
	"log/slog"
	"net/http"
	"time"
)

type HTTPServer struct {
	server          *http.Server
	exchangeHandler *handlers.ExchangeHandler
}

func NewHTTPServer(addr string, exchangeHandler *handlers.ExchangeHandler) *HTTPServer {
	mux := http.NewServeMux()

	httpServer := &HTTPServer{
		exchangeHandler: exchangeHandler,
		server: &http.Server{
			Addr:         addr,
			Handler:      mux,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}

	httpServer.setupRoutes(mux)

	return httpServer
}

func (s *HTTPServer) setupRoutes(mux *http.ServeMux) {
	// Middleware для логирования запросов
	mux.HandleFunc("/api/exchange", s.loggingMiddleware(s.exchangeHandler.ExchangeCurrency))
	mux.HandleFunc("/api/rates", s.loggingMiddleware(s.exchangeHandler.GetExchangeRate))
	mux.HandleFunc("/api/calculate", s.loggingMiddleware(s.exchangeHandler.CalculateExchange))
	mux.HandleFunc("/api/refresh-rates", s.loggingMiddleware(s.exchangeHandler.RefreshRates))
	mux.HandleFunc("/api/cached-rates", s.loggingMiddleware(s.exchangeHandler.GetCachedRates))

	// Health check endpoint
	mux.HandleFunc("/health", s.loggingMiddleware(s.healthCheck))
}

func (s *HTTPServer) loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Добавляем CORS заголовки
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Обрабатываем preflight запросы
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		slog.Info("HTTP request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", duration,
			"remote_addr", r.RemoteAddr,
		)
	}
}

func (s *HTTPServer) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok", "service": "exchange-service"}`))
}

func (s *HTTPServer) Start() error {
	slog.Info("Starting HTTP server", "addr", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	slog.Info("Stopping HTTP server")
	return s.server.Shutdown(ctx)
}
