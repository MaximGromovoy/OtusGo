package interfaces

import "OtusGo/internal/transaction"

// TransactionRepositoryInterface определяет контракт для репозитория транзакций
type TransactionRepositoryInterface interface {
	Add(tx *transaction.Transaction) error
	Get(id int) (*transaction.Transaction, error)
	GetAll() []*transaction.Transaction
	Update(tx *transaction.Transaction) error
}
