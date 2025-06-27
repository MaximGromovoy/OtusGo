package transactionsRepository

const (
	queryInsertTransaction = `
        INSERT INTO transactions (user_id, from_currency, to_currency, from_amount, to_amount, exchange_rate, commission, status, timestamp)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id
    `
	querySelectTransactionByID = `
        SELECT id, user_id, from_currency, to_currency, from_amount, to_amount, exchange_rate, commission, status, timestamp
        FROM transactions WHERE id = $1
    `
	querySelectAllTransactions = `
        SELECT id, user_id, from_currency, to_currency, from_amount, to_amount, exchange_rate, commission, status, timestamp
        FROM transactions
    `
	queryUpdateTransaction = `
        UPDATE transactions
        SET user_id=$1, from_currency=$2, to_currency=$3, from_amount=$4, to_amount=$5, exchange_rate=$6, commission=$7, status=$8, timestamp=$9
        WHERE id=$10
    `
)
