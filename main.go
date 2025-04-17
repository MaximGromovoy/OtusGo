package main

import (
	"OtusGo/internal/model/currency"
	"OtusGo/internal/repository/currencyRepository"
	"OtusGo/internal/service/currencyRepositoryWatcher"
	"OtusGo/internal/service/randomCurrencyService"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	repository := currencyRepository.NewCurrencyRepository()
	watcher := currencyRepositoryWatcher.NewCurrencyRepositoryWatcher(repository)
	randomCurrencyService := randomCurrencyService.NewRandomCurrencyService()

	watcher.StartWatching()

	currencyChannel := make(chan currency.CurrencyInterface)
	randomCurrencyService.StartGenerateCurrencies(currencyChannel, &wg)
	repository.StartListenAndDistribute(currencyChannel, &wg)

	wg.Wait()

	watcher.StopWatching()
}
