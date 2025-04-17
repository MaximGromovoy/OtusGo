package currencyRepository

import (
	"OtusGo/internal/model/currency"
	"sync"
)

func (repo *CurrencyRepository) StartListenAndDistribute(currencyChannel <-chan currency.CurrencyInterface, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for currencyItem := range currencyChannel {
			repo.AddCurrency(currencyItem)
		}
	}()
}

func (repo *CurrencyRepository) AddCurrency(c currency.CurrencyInterface) {
	switch currencyType := c.(type) {
	case *currency.Ruble:
		repo.rubles.mu.Lock()
		defer repo.rubles.mu.Unlock()
		repo.rubles.data = append(repo.rubles.data, currencyType)
	case *currency.Dollar:
		repo.dollars.mu.Lock()
		defer repo.dollars.mu.Unlock()
		repo.dollars.data = append(repo.dollars.data, currencyType)
	case *currency.Euro:
		repo.euros.mu.Lock()
		defer repo.euros.mu.Unlock()
		repo.euros.data = append(repo.euros.data, currencyType)
	case *currency.Lira:
		repo.liras.mu.Lock()
		defer repo.liras.mu.Unlock()
		repo.liras.data = append(repo.liras.data, currencyType)
	}
}
