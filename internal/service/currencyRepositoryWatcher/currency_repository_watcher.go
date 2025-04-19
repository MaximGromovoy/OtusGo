package currencyRepositoryWatcher

import (
	"OtusGo/internal/model/currency"
	"OtusGo/internal/repository/currencyRepository"
	"fmt"
	"time"
)

func StartWatching(repository *currencyRepository.CurrencyRepository, stopChannel <-chan struct{}) {
	go func() {
		fmt.Println("Starting watcher goroutine...")
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()
		currencyCountMap := make(map[string]int)
		for _, cur := range currency.ExistCurrencies {
			currencyCountMap[cur] = 0
		}
		for {
			select {
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
			case <-stopChannel:
				return
			}
		}
	}()
}

func StopWatching(stopChannel chan struct{}) {
	fmt.Println("Stopping watcher goroutine...")
	close(stopChannel)
}
