package transactionRepository

import (
	"OtusGo/internal/converters/transactionConverter"
	"OtusGo/internal/transaction"
	"fmt"
	"sort"
)

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
