package main

import (
	"OtusGo/internal/model/currency"
	"OtusGo/internal/repository/currencyRepository"
	"OtusGo/internal/service/currencyDistributor"
	"OtusGo/internal/service/currencyRepositoryWatcher"
	"OtusGo/internal/service/randomCurrency"
	"context"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	repository := currencyRepository.NewCurrencyRepository()
	currencyChannel := make(chan currency.CurrencyInterface)
	stopChannel := make(chan struct{})

	currencyRepositoryWatcher.StartWatching(repository, stopChannel)
	randomCurrency.StartRandomCurrencyGenerator(currencyChannel)
	currencyDistributor.StartDistributeCurrencyInRepository(currencyChannel, repository)

	<-ctx.Done()

	currencyRepositoryWatcher.StopWatching(stopChannel)
	close(currencyChannel)

}
