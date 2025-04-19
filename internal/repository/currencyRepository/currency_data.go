package currencyRepository

import (
	"OtusGo/internal/model/currency"
	"sync"
)

type currencyData[T currency.CurrencyInterface] struct {
	data []T
	mu   sync.Mutex
}

func (cd *currencyData[T]) Add(item T) {
	cd.mu.Lock()
	defer cd.mu.Unlock()
	cd.data = append(cd.data, item)
}

func (cd *currencyData[T]) GetAll() []T {
	cd.mu.Lock()
	defer cd.mu.Unlock()
	return cd.data
}
