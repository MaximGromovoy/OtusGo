package randomCurrency

import (
	"OtusGo/internal/model/currency"
	"math/rand"
	"time"
)

func StartRandomCurrencyGenerator(currencyChannel chan currency.CurrencyInterface) {
	go func() {
		randomSeed := rand.New(rand.NewSource(1))
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				currencyChannel <- GetRandomCurrency(*randomSeed)
			case <-currencyChannel:
				return
			}
		}
	}()
}
