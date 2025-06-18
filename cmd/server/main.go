package main

import (
	currencyv1 "OtusGo/internal/grpc/pb/api/proto/currency/v1"
	grpcserver "OtusGo/internal/grpc/server"
	"OtusGo/internal/repository/currencyRepository"
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Инициализация репозитория
	storageBaseDir := "."
	repository, err := currencyRepository.NewCurrencyRepository(storageBaseDir)
	if err != nil {
		log.Fatalf("Repository init error: %v", err)
	}

	// Настройка gRPC сервера
	grpcPort := "9090"
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}

	grpcServerInstance := grpc.NewServer()
	currencyGRPCServer := grpcserver.NewCurrencyGRPCServer(repository)
	currencyv1.RegisterCurrencyServiceServer(grpcServerInstance, currencyGRPCServer)

	// Обработка сигналов
	go func() {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
		<-signalChannel

		log.Println("Shutting down gRPC server...")
		grpcServerInstance.GracefulStop()
		cancel()
	}()

	log.Printf("gRPC server is running on port %s", grpcPort)

	if err := grpcServerInstance.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
