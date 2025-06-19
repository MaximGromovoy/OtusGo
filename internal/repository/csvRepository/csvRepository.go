package csvRepository

import (
	"encoding/csv"
	"fmt"
	"os"
)

type CSVRepository struct {
	filePath string
}

func NewCSVRepository(filePath string, headers []string) (*CSVRepository, error) {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, createFile(filePath, headers)
		}
		return nil, fmt.Errorf("failed to open CSV file: %v", err)
	}

	defer file.Close()

	return &CSVRepository{
		filePath: filePath,
	}, nil
}

func (r *CSVRepository) GetRecords() ([][]string, error) {
	file, err := os.Open(r.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %v", err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %v", err)
	}

	records = records[1:]

	return records, nil
}

func (r *CSVRepository) AppendRecord(record []string) error {
	file, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open CSV file for append: %v", err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return writer.Write(record)
}

func createFile(filePath string, headers []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return writer.Write(headers)
}
