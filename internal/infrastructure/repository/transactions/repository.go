package transactionsRepository

import (
	"OtusGo/internal/application/models/transaction"
	postgresRepository "OtusGo/internal/infrastructure/repository/postgres"
	"context"
	"fmt"
	"time"
)

type TransactionsPostgresRepository struct {
	db *postgresRepository.PostgresDB
}

func NewTransactionsPostgresRepository(db *postgresRepository.PostgresDB) (*TransactionsPostgresRepository, error) {
	return &TransactionsPostgresRepository{db: db}, nil
}

func (r *TransactionsPostgresRepository) Add(tx *transaction.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.db.QueryRow(ctx, queryInsertTransaction,
		tx.UserID,
		tx.FromCurrency,
		tx.ToCurrency,
		tx.FromAmount,
		tx.ToAmount,
		tx.ExchangeRate,
		tx.Commission,
		tx.Status,
		tx.Timestamp,
	).Scan(&tx.ID)
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %v", err)
	}
	return nil
}

func (r *TransactionsPostgresRepository) Get(id int) (*transaction.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var tx transaction.Transaction
	err := r.db.QueryRow(ctx, querySelectTransactionByID, id).Scan(
		&tx.ID,
		&tx.UserID,
		&tx.FromCurrency,
		&tx.ToCurrency,
		&tx.FromAmount,
		&tx.ToAmount,
		&tx.ExchangeRate,
		&tx.Commission,
		&tx.Status,
		&tx.Timestamp,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find transaction: %v", err)
	}
	return &tx, nil
}

func (r *TransactionsPostgresRepository) GetAll() []*transaction.Transaction {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := r.db.Query(ctx, querySelectAllTransactions)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var transactions []*transaction.Transaction
	for rows.Next() {
		var tx transaction.Transaction
		if err := rows.Scan(
			&tx.ID,
			&tx.UserID,
			&tx.FromCurrency,
			&tx.ToCurrency,
			&tx.FromAmount,
			&tx.ToAmount,
			&tx.ExchangeRate,
			&tx.Commission,
			&tx.Status,
			&tx.Timestamp,
		); err != nil {
			continue
		}
		transactions = append(transactions, &tx)
	}
	return transactions
}

func (r *TransactionsPostgresRepository) Update(tx *transaction.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd, err := r.db.Pool.Exec(ctx, queryUpdateTransaction,
		tx.UserID,
		tx.FromCurrency,
		tx.ToCurrency,
		tx.FromAmount,
		tx.ToAmount,
		tx.ExchangeRate,
		tx.Commission,
		tx.Status,
		tx.Timestamp,
		tx.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %v", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("transaction with ID %d not found", tx.ID)
	}
	return nil
}
