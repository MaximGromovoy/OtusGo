package currencyRepositoryWatcher

import (
	"OtusGo/internal/model/currency"
	"OtusGo/internal/repository/currencyRepository"
	"context"
	"fmt"
	"time"
)

func StartWatching(repository *currencyRepository.CurrencyRepository, ctx context.Context) {
	fmt.Println("Starting watcher goroutine...")

	// Инициализируем карту для хранения текущего количества записей по каждой валюте
	currencyCountMap := make(map[string]int)

	// Заполняем начальными значениями из репозитория
	for _, curName := range currency.ExistCurrencies {
		data := repository.GetAll(curName)
		currencyCountMap[curName] = len(data)
		fmt.Printf("Initial count for %s: %d records\n", curName, len(data))
	}

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Watcher stopped.")
			return
		case <-ticker.C:
			for curName, curSize := range currencyCountMap {
				var data = repository.GetAll(curName)
				var currentSize = len(data)
				if curSize < currentSize {
					currencyCountMap[curName] = currentSize
					for i := curSize; i < currentSize; i++ {
						item := data[i]
						fmt.Printf("New currency %-10s | Code: %-5s | Value: %10.2f\n",
							item.GetName(), item.GetCode(), item.GetValue())
					}
				}
			}
		}
	}
}
