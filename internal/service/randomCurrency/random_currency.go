package randomCurrency

import (
	"OtusGo/internal/model/currency"
	"math/rand"
)

func GetRandomCurrency(randomSeed rand.Rand) currency.CurrencyInterface {
	index := randomSeed.Intn(len(currency.SupportedCurrencies))
	currencyName := currency.SupportedCurrencies[index]

	return currency.NewCurrency(currencyName, randomSeed.Float64()*100)
}
