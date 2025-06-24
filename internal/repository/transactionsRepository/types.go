package transactionRepository

import (
	"OtusGo/internal/converters/transactionConverter"
	"OtusGo/internal/interfaces"
	"OtusGo/internal/repository/csvRepository"
	"OtusGo/internal/transaction"
	"log/slog"
	"path/filepath"
	"sync"
)

var _ interfaces.TransactionRepositoryInterface = (*TransactionRepository)(nil)

const subDir = "transactionsStorage"

type TransactionRepository struct {
	transactions  map[int]*transaction.Transaction
	csvRepository *csvRepository.CSVRepository
	lastID        int
	mutex         sync.RWMutex
}

// NewTransactionRepository создает новый репозиторий транзакций
func NewTransactionRepository(baseDir string) (*TransactionRepository, error) {
	var headers = transactionConverter.GetHeaders()

	csvRepository, err := csvRepository.NewCSVRepository(
		filepath.Join(baseDir, subDir, "transactions.csv"), headers)

	if err != nil {
		return nil, err
	}

	records, err := csvRepository.GetRecords()
	if err != nil {
		return nil, err
	}

	var transactionModels = make(map[int]*transaction.Transaction)
	var lastID int = 0

	for i, record := range records {
		transactionModel, err := transactionConverter.ToModel(record)

		if err != nil {
			slog.Warn("Failed to convert transaction record", "error", err, "record_index", i)
			continue
		}

		transactionModels[i] = transactionModel

		if transactionModel.ID > lastID {
			lastID = transactionModel.ID
		}
	}

	repo := &TransactionRepository{
		transactions:  transactionModels,
		csvRepository: csvRepository,
		mutex:         sync.RWMutex{},
		lastID:        lastID,
	}

	return repo, nil
}
