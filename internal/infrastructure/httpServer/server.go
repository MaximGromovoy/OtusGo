package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"OtusGo/internal/infrastructure/httpServer/handlers"
)

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(addr string, handler *handlers.ExchangeHandler) *HTTPServer {
	mux := http.NewServeMux()

	// Регистрируем обработчики
	mux.HandleFunc("/api/exchange", handler.Exchange)
	mux.HandleFunc("/api/rate", handler.GetRate)

	return &HTTPServer{
		server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

func (s *HTTPServer) Start() error {
	// Запускаем сервер в отдельной горутине
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Ожидаем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}

func (s *HTTPServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}
