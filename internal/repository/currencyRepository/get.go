package currencyRepository

import (
	"OtusGo/internal/model/currency"
)

func (repo *CurrencyRepository) GetRubles() []*currency.Ruble {
	repo.rubles.mu.Lock()
	defer repo.rubles.mu.Unlock()
	return repo.rubles.data
}

func (repo *CurrencyRepository) GetDollars() []*currency.Dollar {
	repo.dollars.mu.Lock()
	defer repo.dollars.mu.Unlock()
	return repo.dollars.data
}

func (repo *CurrencyRepository) GetEuros() []*currency.Euro {
	repo.euros.mu.Lock()
	defer repo.euros.mu.Unlock()
	return repo.euros.data
}

func (repo *CurrencyRepository) GetLiras() []*currency.Lira {
	repo.liras.mu.Lock()
	defer repo.liras.mu.Unlock()
	return repo.liras.data
}
