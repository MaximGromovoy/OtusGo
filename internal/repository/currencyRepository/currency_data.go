package currencyRepository

import (
	"OtusGo/internal/model/currency"
	"errors"
	"sync"
)

// CurrencyCollection представляет интерфейс для работы с коллекцией валют
type CurrencyCollection interface {
	Add(item currency.CurrencyInterface) error
	GetAll() []currency.CurrencyInterface
	Update(id int, item currency.CurrencyInterface) error
	Delete(id int) error
	GetByID(id int) (currency.CurrencyInterface, bool)
}

// currencyData представляет собой хранилище данных для конкретного типа валюты.
type currencyData struct {
	items    []currency.CurrencyInterface // срез элементов валюты
	filePath string                       // путь к CSV-файлу для сохранения данных
	mu       sync.RWMutex                 // привный мьютекс для синхронизации доступа к данным
}

// NewCurrencyData создает новую коллекцию валют
func NewCurrencyData(filePath string, initialItems []currency.CurrencyInterface) *currencyData {
	return &currencyData{
		items:    initialItems,
		filePath: filePath,
	}
}

// Add добавляет новый элемент в коллекцию и сохраняет его в файл
func (c *currencyData) Add(item currency.CurrencyInterface) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Добавляем новый элемент в срез
	c.items = append(c.items, item)

	// Сохраняем данные в файл
	return appendCurrencyToFile(c.filePath, item)
}

// GetAll возвращает копию всех элементов из хранилища.
func (c *currencyData) GetAll() []currency.CurrencyInterface {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Создаем новый срез того же размера
	result := make([]currency.CurrencyInterface, len(c.items))

	// Копируем элементы
	copy(result, c.items)

	return result
}

// Update обновляет элемент по его ID
func (c *currencyData) Update(id int, item currency.CurrencyInterface) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	found := false
	for i, curr := range c.items {
		if curr.GetID() == id {
			c.items[i] = item
			found = true
			break
		}
	}

	if !found {
		return errors.New("currency not found")
	}

	// Перезаписываем все в файл
	return rewriteCurrenciesInFile(c.filePath, c.items)
}

// Delete удаляет элемент по его ID
func (c *currencyData) Delete(id int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, curr := range c.items {
		if curr.GetID() == id {
			// Удаляем элемент из среза
			c.items = append(c.items[:i], c.items[i+1:]...)

			// Перезаписываем все в файл
			return rewriteCurrenciesInFile(c.filePath, c.items)
		}
	}

	return errors.New("currency not found")
}

// GetByID возвращает элемент по его ID
func (c *currencyData) GetByID(id int) (currency.CurrencyInterface, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, curr := range c.items {
		if curr.GetID() == id {
			return curr, true
		}
	}

	return nil, false
}

// GetFilePath возвращает путь к файлу хранилища
func (c *currencyData) GetFilePath() string {
	return c.filePath
}
