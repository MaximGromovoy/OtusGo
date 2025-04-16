package currencyRepository

import (
	"OtusGo/internal/model/currency"
)

type CurrencyRepository struct {
	rubles  []*currency.Ruble
	dollars []*currency.Dollar
	euros   []*currency.Euro
	liras   []*currency.Lira
}

// NewCurrencyRepository создает новый экземпляр репозитория
func NewCurrencyRepository() *CurrencyRepository {
	return &CurrencyRepository{
		rubles:  make([]*currency.Ruble, 0),
		dollars: make([]*currency.Dollar, 0),
		euros:   make([]*currency.Euro, 0),
		liras:   make([]*currency.Lira, 0),
	}
}

func (repo *CurrencyRepository) AddCurrencies(currencies []currency.CurrencyInterface) {
	for _, currencyItem := range currencies {
		repo.AddCurrency(currencyItem)
	}
}

func (repo *CurrencyRepository) AddCurrency(c currency.CurrencyInterface) {
	switch currencyType := c.(type) {
	case *currency.Ruble:
		repo.rubles = append(repo.rubles, currencyType)
	case *currency.Dollar:
		repo.dollars = append(repo.dollars, currencyType)
	case *currency.Euro:
		repo.euros = append(repo.euros, currencyType)
	case *currency.Lira:
		repo.liras = append(repo.liras, currencyType)
	}
}

func (repo *CurrencyRepository) GetRubles() []*currency.Ruble {
	return repo.rubles
}

func (repo *CurrencyRepository) GetDollars() []*currency.Dollar {
	return repo.dollars
}

func (repo *CurrencyRepository) GetEuros() []*currency.Euro {
	return repo.euros
}

func (repo *CurrencyRepository) GetLiras() []*currency.Lira {
	return repo.liras
}
