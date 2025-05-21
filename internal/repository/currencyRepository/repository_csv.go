package currencyRepository

import (
	"OtusGo/internal/model/currency"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

func fetchCurrenciesFromFile(filePath string, currencyName string) (int, []currency.CurrencyInterface, error) {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, make([]currency.CurrencyInterface, 0), nil
		}
		return 0, nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var items []currency.CurrencyInterface
	maxID := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return maxID, nil, fmt.Errorf("failed to read record from %s: %w", filePath, err)
		}

		if len(record) < 3 {
			fmt.Printf("Warning: skipping invalid record in %s: %v\n", filePath, record)
			continue
		}

		id := 0
		if parsedID, err := strconv.Atoi(record[0]); err == nil {
			id = parsedID
			if id > maxID {
				maxID = id
			}
		}

		nameFromFile := record[1]
		if nameFromFile != currencyName {
			fmt.Printf("Warning: skipping record with mismatched currency name in %s. Expected: %s, Got: %s\n",
				filePath, currencyName, nameFromFile)
			continue
		}

		value, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			fmt.Printf("Warning: skipping record with invalid value in %s: %v\n", filePath, record)
			continue
		}

		item := currency.NewCurrencyWithID(nameFromFile, value, id)
		items = append(items, item)
	}

	return maxID, items, nil
}

func rewriteCurrenciesInFile(filePath string, collection []currency.CurrencyInterface) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s for overwriting: %w", filePath, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, curr := range collection {
		record := getRecord(curr)
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record to %s: %w", filePath, err)
		}
	}

	return nil
}

func appendCurrencyToFile(filePath string, savedCurrency currency.CurrencyInterface) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %s for appending: %w", filePath, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := getRecord(savedCurrency)
	if err := writer.Write(record); err != nil {
		return fmt.Errorf("failed to write record to %s: %w", filePath, err)
	}

	return nil
}

func getRecord(currency currency.CurrencyInterface) []string {
	return []string{
		fmt.Sprintf("%d", currency.GetID()),
		currency.GetName(),
		fmt.Sprintf("%.2f", currency.GetValue()),
	}
}
