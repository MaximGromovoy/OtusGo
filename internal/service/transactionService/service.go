package transactionService

import (
	transactionRepository "OtusGo/internal/repository/transactionsRepository"
	transaction "OtusGo/internal/transaction"
)

type TransactionService struct {
	transactionRepository *transactionRepository.TransactionRepository
}

var depositCommission = 0.0
var withdrawCommission = 0.025

//var exchangeCommission = 0.05

func NewTransactionService(repo *transactionRepository.TransactionRepository) *TransactionService {
	return &TransactionService{
		transactionRepository: repo,
	}
}

func (service *TransactionService) Deposit(
	userID int, currency string, amount float64) {
	var exchangeRate = 1.0

	var tx = transaction.NewDepositTransaction(
		userID, currency,
		currency, amount,
		amount, exchangeRate, depositCommission)

	tx.Complete()
	err := service.transactionRepository.Add(tx)

	if err != nil {
		tx.Fail()
		err = service.transactionRepository.Add(tx)

		if err != nil {
			panic("Failed deposit transaction: " + err.Error())
		}
	}
}

func (service *TransactionService) Withdraw(
	userID int, currency string, amount float64) {
	var exchangeRate = 1.0

	var tx = transaction.NewWithdrawTransaction(
		userID, currency,
		currency, amount,
		amount, exchangeRate, withdrawCommission)

	tx.Complete()
	err := service.transactionRepository.Add(tx)

	if err != nil {
		tx.Fail()
		err = service.transactionRepository.Add(tx)

		if err != nil {
			panic("Failed withdraw transaction: " + err.Error())
		}
	}
}
