package randomCurrency

import (
	"OtusGo/internal/model/currency"
	"math/rand"
)

func GetRandomCurrency(randomSeed rand.Rand) currency.CurrencyInterface {
	index := randomSeed.Intn(len(currency.ExistCurrencies))
	currencyName := currency.ExistCurrencies[index]
	switch currencyName {
	case currency.RubleCurrency:
		return currency.NewRuble(randomSeed.Float64() * 100)
	case currency.DollarCurrency:
		return currency.NewDollar(randomSeed.Float64() * 100)
	case currency.EuroCurrency:
		return currency.NewEuro(randomSeed.Float64() * 100)
	case currency.LiraCurrency:
		return currency.NewLira(randomSeed.Float64() * 100)
	}
	return nil
}
