package transactionRepository

import (
	"OtusGo/internal/converters/transactionConverter"
	"OtusGo/internal/repository/csvRepository"
	"OtusGo/internal/transaction"
	"fmt"
	"log/slog"
	"path/filepath"
	"sort"
	"sync"
)

type TransactionRepository struct {
	transactions  map[int]*transaction.Transaction
	csvRepository *csvRepository.CSVRepository
	lastID        int
	mutex         sync.RWMutex
}

const subDir = "transactionsStorage"

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

// GetByID возвращает транзакцию по ID
func (r *TransactionRepository) Get(id int) (*transaction.Transaction, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tx, exists := r.transactions[id]
	if !exists {
		return nil, fmt.Errorf("transaction with ID %d not found", id)
	}

	return tx, nil
}

// GetAll возвращает все транзакции
func (r *TransactionRepository) GetAll() []*transaction.Transaction {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var allTransactions []*transaction.Transaction
	for _, tx := range r.transactions {
		allTransactions = append(allTransactions, tx)
	}

	sort.Slice(allTransactions, func(i, j int) bool {
		return allTransactions[i].Timestamp.After(allTransactions[j].Timestamp)
	})

	return allTransactions
}

// Add добавляет новую транзакцию
func (r *TransactionRepository) Add(tx *transaction.Transaction) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	tx.ID = r.lastID + 1
	r.transactions[tx.ID] = tx

	record, err := transactionConverter.ToRecord(tx)

	if err != nil {
		return err
	}

	err = r.csvRepository.AppendRecord(record)

	if err != nil {
		r.transactions[tx.ID] = nil
	}

	return err
}
