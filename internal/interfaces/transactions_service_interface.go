package interfaces

import "OtusGo/internal/transaction"

type TransactionsServiceInterface interface {
	// SaveTransaction сохраняет транзакцию в сервисе
	SaveTransaction(transaction *transaction.Transaction) error
}
