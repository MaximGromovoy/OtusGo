package main

import (
	"OtusGo/internal/model/currency"
	"OtusGo/internal/repository/currencyRepository"
	"OtusGo/internal/service/randomCurrencyService"
	"OtusGo/internal/service/randomCurrencyService/currencyRepositoryWatcher"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	repository := currencyRepository.NewCurrencyRepository()
	currencyRepositoryWatcher.NewCurrencyRepositoryWatcher(repository)

	randomCurrencyService := randomCurrencyService.NewRandomCurrencyService()

	currencyChannel := make(chan currency.CurrencyInterface)

	wg.Add(1)
	go func() {
		defer wg.Done()
		randomCurrencyService.GenerateCurrencies(currencyChannel)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		repository.ListenAndDistribute(currencyChannel)
	}()

	wg.Wait()
}
