package repository

import (
	"OtusGo/internal/model/currency"
)

func DistributeCurrencies(currencies []currency.CurrencyInterface) (rubles []*currency.Ruble,
	dollars []*currency.Dollar, euros []*currency.Euro, liras []*currency.Lira) {
	for _, currencyItem := range currencies {
		switch currencyType := currencyItem.(type) {
		case *currency.Ruble:
			rubles = append(rubles, currencyType)
		case *currency.Dollar:
			dollars = append(dollars, currencyType)
		case *currency.Euro:
			euros = append(euros, currencyType)
		case *currency.Lira:
			liras = append(liras, currencyType)
		}
	}
	return
}
