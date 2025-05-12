package currencyRepository

import (
	"OtusGo/internal/model/currency"
	"encoding/csv"
	"fmt"
	"os"
	"sync"
)

// currencyData представляет собой потокобезопасное хранилище для определенного типа валюты.
// Оно управляет данными в памяти и обеспечивает их сохранение в CSV-файл.
type currencyData[T currency.CurrencyInterface] struct {
	data     []T
	mu       sync.Mutex
	filePath string
}

// Add добавляет новый элемент валюты в хранилище.
// Элемент добавляется как в срез в памяти, так и в CSV-файл.
func (cd *currencyData[T]) Add(item T) error {
	cd.mu.Lock()
	defer cd.mu.Unlock()

	cd.data = append(cd.data, item)

	file, err := os.OpenFile(cd.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл %s для дозаписи: %w", cd.filePath, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	record := []string{item.GetName(), item.GetCode(), fmt.Sprintf("%.2f", item.GetValue())}
	if err := writer.Write(record); err != nil {
		writer.Flush()
		return fmt.Errorf("не удалось записать запись в %s: %w", cd.filePath, err)
	}

	if err := writer.Error(); err != nil {
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("не удалось очистить буфер для %s: %w", cd.filePath, err)
	}

	return nil
}

// GetAll возвращает копию всех элементов валюты, хранящихся в данный момент.
// Возвращает новый срез, чтобы избежать модификации внутренних данных извне.
func (cd *currencyData[T]) GetAll() []T {
	cd.mu.Lock()
	defer cd.mu.Unlock()
	dataCopy := make([]T, len(cd.data))
	copy(dataCopy, cd.data)
	return dataCopy
}
