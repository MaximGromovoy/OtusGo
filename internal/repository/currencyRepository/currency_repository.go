package currencyRepository

import (
	"OtusGo/internal/model/currency"
)

type CurrencyRepository struct {
	data map[string]*currencyData[currency.CurrencyInterface]
}

func NewCurrencyRepository() *CurrencyRepository {
	data := make(map[string]*currencyData[currency.CurrencyInterface])

	for _, currencyType := range currency.ExistCurrencies {
		if _, exists := data[currencyType]; !exists {
			data[currencyType] = &currencyData[currency.CurrencyInterface]{data: make([]currency.CurrencyInterface, 0)}
		}
	}

	return &CurrencyRepository{data: data}
}

func (repo *CurrencyRepository) Add(item currency.CurrencyInterface) {
	if collection, exists := repo.data[item.GetName()]; exists {
		collection.Add(item)
	}
}

func (repo *CurrencyRepository) GetAll(currencyType string) []currency.CurrencyInterface {
	if collection, exists := repo.data[currencyType]; exists {
		return collection.GetAll()
	}
	return nil
}
