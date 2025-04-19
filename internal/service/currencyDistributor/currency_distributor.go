package currencyDistributor

import (
	"OtusGo/internal/model/currency"
	"OtusGo/internal/repository/currencyRepository"
	"context"
	"fmt"
)

func StartDistributeCurrencyInRepository(currencyChannel <-chan currency.CurrencyInterface,
	repository *currencyRepository.CurrencyRepository, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Currency distributor stopped.")
			return
		case currency, ok := <-currencyChannel:
			if !ok {
				fmt.Println("Currency channel closed. Distributor exiting.")
				return
			}
			repository.Add(currency)
			fmt.Printf("Distributed currency: %s\n", currency.GetName())
		}
	}
}
