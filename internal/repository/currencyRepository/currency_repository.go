package currencyRepository

import (
	"OtusGo/internal/model/currency"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

// storageSubDir определяет имя поддиректории для хранения CSV-файлов с данными о валютах.
const storageSubDir = "currency_storage"

// CurrencyRepository представляет собой репозиторий для управления данными о различных валютах.
// Он организует хранение данных по типам валют, используя для каждого типа отдельный
// экземпляр currencyData. Данные сохраняются в CSV-файлах в указанной директории.
type CurrencyRepository struct {
	data       map[string]*currencyData[currency.CurrencyInterface] // Карта для хранения данных по типам валют. Ключ - имя валюты.
	storageDir string                                               // Путь к директории, где хранятся CSV-файлы.
}

// NewCurrencyRepository создает и инициализирует новый экземпляр CurrencyRepository.
// baseDir - это базовая директория, внутри которой будет создана поддиректория storageSubDir
// для хранения файлов данных.
func NewCurrencyRepository(baseDir string) (*CurrencyRepository, error) {
	repoStorageDir := filepath.Join(baseDir, storageSubDir)
	if err := os.MkdirAll(repoStorageDir, 0755); err != nil {
		return nil, fmt.Errorf("не удалось создать директорию для хранения %s: %w", repoStorageDir, err)
	}

	dataMap := make(map[string]*currencyData[currency.CurrencyInterface])

	for _, currencyName := range currency.ExistCurrencies {
		filePath := filepath.Join(repoStorageDir, currencyName+".csv")

		loadedItems, err := loadItemsFromFile(filePath, currencyName)
		if err != nil {
			// Логируем предупреждение, но не прерываем работу,
			// чтобы приложение могло работать с пустым списком для этой валюты.
			fmt.Printf("Предупреждение: не удалось загрузить данные для %s из %s: %v. Будет использован пустой список.\n", currencyName, filePath, err)
			loadedItems = make([]currency.CurrencyInterface, 0) // Используем пустой срез в случае ошибки
		}

		dataMap[currencyName] = &currencyData[currency.CurrencyInterface]{
			data:     loadedItems,
			filePath: filePath,
		}
	}

	return &CurrencyRepository{
		data:       dataMap,
		storageDir: repoStorageDir,
	}, nil
}

// loadItemsFromFile загружает элементы валюты из указанного CSV-файла.
// filePath - путь к CSV-файлу.
// currencyNameTarget - ожидаемое имя валюты, записи с другим именем будут проигнорированы.
func loadItemsFromFile(filePath string, currencyNameTarget string) ([]currency.CurrencyInterface, error) {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return make([]currency.CurrencyInterface, 0), nil
		}
		return nil, fmt.Errorf("не удалось открыть файл %s: %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var items []currency.CurrencyInterface

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Не удалось прочитать запись из %s: %w", filePath, err)
		}

		if len(record) != 3 {
			fmt.Printf("Предупреждение: пропуск некорректной записи в %s: %v\n", filePath, record)
			continue
		}

		nameFromFile := record[0]
		if nameFromFile != currencyNameTarget {
			fmt.Printf("Предупреждение: пропуск записи с несоответствующим именем валюты в %s. Ожидалось: %s, Получено: %s\n", filePath, currencyNameTarget, nameFromFile)
			continue
		}

		value, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			fmt.Printf("Предупреждение: пропуск записи с неверным значением в %s: %v\n", filePath, record)
			continue
		}

		var item currency.CurrencyInterface
		switch currencyNameTarget {
		case currency.RubleCurrency:
			item = currency.NewRuble(value)
		case currency.DollarCurrency:
			item = currency.NewDollar(value)
		case currency.EuroCurrency:
			item = currency.NewEuro(value)
		case currency.LiraCurrency:
			item = currency.NewLira(value)
		default:
			fmt.Printf("Предупреждение: неизвестный тип валюты %s при загрузке из файла %s\n", currencyNameTarget, filePath)
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

// Add добавляет новый элемент валюты в соответствующее хранилище (currencyData).
// Тип валюты определяется методом GetName() переданного элемента.
// Если для данного типа валюты хранилище не инициализировано, возвращается ошибка.
// В противном случае, элемент добавляется в память и сохраняется в соответствующий CSV-файл.
func (repo *CurrencyRepository) Add(item currency.CurrencyInterface) error {
	currencyName := item.GetName()
	collection, exists := repo.data[currencyName]
	if !exists {
		return fmt.Errorf("Тип валюты %s не инициализирован в репозитории", currencyName)
	}
	return collection.Add(item)
}

// GetAll возвращает срез всех элементов для указанного типа валюты.
// currencyType - строка, представляющая тип валюты
// Если данные для указанного типа валюты существуют, возвращается копия среза этих данных.
// Если тип валюты не найден в репозитории, возвращается nil.
func (repo *CurrencyRepository) GetAll(currencyType string) []currency.CurrencyInterface {
	if collection, exists := repo.data[currencyType]; exists {
		return collection.GetAll()
	}
	return nil
}
