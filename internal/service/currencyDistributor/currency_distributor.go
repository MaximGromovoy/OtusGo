package currencyDistributor

import (
	"OtusGo/internal/model/currency"
	"OtusGo/internal/repository/currencyRepository"
)

func StartDistributeCurrencyInRepository(currencyChannel <-chan currency.CurrencyInterface,
	repository *currencyRepository.CurrencyRepository) {
	go func() {
		for item := range currencyChannel {
			repository.Add(item)
		}
	}()
}
