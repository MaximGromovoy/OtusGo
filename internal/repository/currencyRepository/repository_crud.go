package currencyRepository

import (
	"OtusGo/internal/model/currency"
	"errors"
)

// Add добавляет новую валюту в репозиторий
func (repo *CurrencyRepository) Add(item currency.CurrencyInterface) error {
	currName := item.GetName()
	collection, exists := repo.data[currName]
	if !exists {
		return errors.New("currency not initialized")
	}

	return collection.Add(item)
}

// GetAll возвращает все валюты указанного типа
func (repo *CurrencyRepository) GetAll(currName string) []currency.CurrencyInterface {
	collection, exists := repo.data[currName]
	if !exists {
		return []currency.CurrencyInterface{}
	}

	return collection.GetAll()
}

// Update обновляет валюту указанного типа по ID
func (repo *CurrencyRepository) Update(currName string, id int, item currency.CurrencyInterface) error {
	collection, exists := repo.data[currName]
	if !exists {
		return errors.New("currency not initialized")
	}

	return collection.Update(id, item)
}

// DeleteByTypeAndID удаляет валюту указанного типа по ID
func (repo *CurrencyRepository) DeleteByTypeAndID(currName string, id int) error {
	collection, exists := repo.data[currName]
	if !exists {
		return errors.New("currency not initialized")
	}

	return collection.Delete(id)
}

// GetByTypeAndID возвращает валюту указанного типа по ID
func (repo *CurrencyRepository) GetByTypeAndID(currName string, id int) (currency.CurrencyInterface, error) {
	collection, exists := repo.data[currName]
	if !exists {
		return nil, errors.New("currency not initialized")
	}

	item, found := collection.GetByID(id)
	if !found {
		return nil, errors.New("currency not found")
	}

	return item, nil
}
