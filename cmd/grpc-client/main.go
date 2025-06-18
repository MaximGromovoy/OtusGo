package main

import (
	currencyv1 "OtusGo/internal/grpc/pb/api/proto/currency/v1"
	"OtusGo/internal/model/currency"
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Подключение к gRPC серверу
	conn, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := currencyv1.NewCurrencyServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Тестирование различных операций
	testCRUDOperations(client, ctx)
}

func testCRUDOperations(client currencyv1.CurrencyServiceClient, ctx context.Context) {
	fmt.Println("=== Testing gRPC Currency Service ===")

	// 1. Создание валют
	fmt.Println("\n1. Creating currencies...")

	currencies := []struct {
		Type  string
		Value float64
	}{
		{currency.DollarCurrency, 100.50},
		{currency.EuroCurrency, 85.25},
		{currency.RubleCurrency, 7500.00},
		{currency.LiraCurrency, 1850.75},
	}

	var createdIDs []int32

	for _, curr := range currencies {
		resp, err := client.CreateCurrency(ctx, &currencyv1.CreateCurrencyRequest{
			CurrencyType: curr.Type,
			Value:        curr.Value,
		})
		if err != nil {
			log.Printf("Failed to create %s: %v", curr.Type, err)
			continue
		}

		fmt.Printf("Created %s: ID=%d, Value=%.2f, Code=%s\n",
			resp.Currency.Name, resp.Currency.Id, resp.Currency.Value, resp.Currency.Code)
		createdIDs = append(createdIDs, resp.Currency.Id)
	}

	// 2. Получение всех валют по типам
	fmt.Println("\n2. Getting all currencies by type...")

	for _, currType := range []string{currency.DollarCurrency, currency.EuroCurrency, currency.RubleCurrency, currency.LiraCurrency} {
		resp, err := client.GetCurrencies(ctx, &currencyv1.GetCurrenciesRequest{
			CurrencyType: currType,
		})
		if err != nil {
			log.Printf("Failed to get %s currencies: %v", currType, err)
			continue
		}

		fmt.Printf("%s currencies (%d items):\n", currType, len(resp.Currencies))
		for _, curr := range resp.Currencies {
			fmt.Printf("  ID=%d, Value=%.2f, Code=%s\n", curr.Id, curr.Value, curr.Code)
		}
	}

	fmt.Println("\n=== Testing completed ===")
}
