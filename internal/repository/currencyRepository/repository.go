package currencyRepository

import (
	"OtusGo/internal/model/currency"
	"fmt"
	"os"
	"path/filepath"
)

// storageSubDir определяет имя поддиректории для хранения CSV-файлов с данными о валютах.
const storageSubDir = "currency_storage"

// CurrencyRepository представляет собой репозиторий для управления данными о различных валютах.
type CurrencyRepository struct {
	data       map[string]*currencyData // Карта для хранения данных по типам валют. Ключ - имя валюты.
	storageDir string                   // Путь к директории, где хранятся CSV-файлы.
}

// NewCurrencyRepository создает и инициализирует новый экземпляр CurrencyRepository.
func NewCurrencyRepository(baseDir string) (*CurrencyRepository, error) {
	repoStorageDir := filepath.Join(baseDir, storageSubDir)
	if err := os.MkdirAll(repoStorageDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory %s: %w", repoStorageDir, err)
	}

	dataMap := make(map[string]*currencyData)
	maxIDs := make(map[string]int)

	for _, currencyName := range currency.SupportedCurrencies {
		filePath := filepath.Join(repoStorageDir, currencyName+".csv")

		maxID, loadedItems, err := fetchCurrenciesFromFile(filePath, currencyName)
		if err != nil {
			fmt.Printf("Warning: failed to load data for %s from %s: %v. An empty list will be used.\n", currencyName, filePath, err)
			loadedItems = make([]currency.CurrencyInterface, 0)
		}

		maxIDs[currencyName] = maxID

		dataMap[currencyName] = NewCurrencyData(filePath, loadedItems)
	}

	currency.InitializeCounters(maxIDs)

	return &CurrencyRepository{
		data:       dataMap,
		storageDir: repoStorageDir,
	}, nil
}
