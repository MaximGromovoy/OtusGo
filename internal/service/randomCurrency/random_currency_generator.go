package randomCurrency

import (
	"OtusGo/internal/model/currency"
	"context"
	"fmt"
	"math/rand"
	"time"
)

func StartRandomCurrencyGenerator(currencyChannel chan currency.CurrencyInterface, ctx context.Context) {
	randomSeed := rand.New(rand.NewSource(1))
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Currency generator stopped.")
			return
		case <-ticker.C:
			currencyChannel <- GetRandomCurrency(*randomSeed)
		}
	}
}
