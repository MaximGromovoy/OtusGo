package currencyRepository

import (
	"OtusGo/internal/model/currency"
	"sync"
)

type CurrencyRepository struct {
	rubles struct {
		data []*currency.Ruble
		mu   sync.Mutex
	}
	dollars struct {
		data []*currency.Dollar
		mu   sync.Mutex
	}
	euros struct {
		data []*currency.Euro
		mu   sync.Mutex
	}
	liras struct {
		data []*currency.Lira
		mu   sync.Mutex
	}
}

func NewCurrencyRepository() *CurrencyRepository {
	return &CurrencyRepository{
		rubles: struct {
			data []*currency.Ruble
			mu   sync.Mutex
		}{
			data: make([]*currency.Ruble, 0),
		},
		dollars: struct {
			data []*currency.Dollar
			mu   sync.Mutex
		}{
			data: make([]*currency.Dollar, 0),
		},
		euros: struct {
			data []*currency.Euro
			mu   sync.Mutex
		}{
			data: make([]*currency.Euro, 0),
		},
		liras: struct {
			data []*currency.Lira
			mu   sync.Mutex
		}{
			data: make([]*currency.Lira, 0),
		},
	}
}
