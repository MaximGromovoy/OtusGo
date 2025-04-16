package service

import (
	"OtusGo/internal/model/currency"
	"math/rand"
	"time"
)

func GetRandomCurrencies() (currencies []currency.CurrencyInterface) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 20; i++ {
		localCurrencies := []currency.CurrencyInterface{
			currency.NewRuble(r.Float64() * 100),
			currency.NewDollar(r.Float64() * 100),
			currency.NewEuro(r.Float64() * 100),
			currency.NewLira(r.Float64() * 100),
		}
		currencies = append(currencies, localCurrencies[r.Intn(len(localCurrencies))])
	}
	return
}
