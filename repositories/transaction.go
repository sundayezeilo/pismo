package repositories

import (
	"context"
	"database/sql"

	"github.com/sundayezeilo/pismo/models"
)

type TxnRepository interface {
	CreateTransaction(context.Context, *CreateTransactionParams) (models.Transaction, error)
	GetTransactionByID(context.Context, int) (*models.Transaction, error)
}

type txnRepository struct {
	db *sql.DB
}

func NewTxnRepository(db *sql.DB) TxnRepository {
	return &txnRepository{db}
}

type CreateTransactionParams struct {
	AccountID int     `json:"account_id"`
	OpTypeID  int     `json:"operation_type_id"`
	Amount    float64 `json:"amount"`
}

func (r *txnRepository) CreateTransaction(ctx context.Context, txn *CreateTransactionParams) (models.Transaction, error) {
	query := `
		WITH account_balance AS (
				SELECT balance
				FROM accounts
				WHERE id = $1
		),
		inserted_transaction AS (
				INSERT INTO transactions (account_id, amount, operation_type_id, balance_before, balance_after)
				VALUES ($1, $2, $3, (SELECT balance FROM account_balance), (SELECT balance FROM account_balance) + $2)
				RETURNING id AS transaction_id, account_id, operation_type_id, amount, event_date, balance_before, balance_after, created_at, updated_at
		),
		updated_account AS (
				UPDATE accounts
				SET balance = balance + $2
				WHERE id = $1
				RETURNING id AS account_id, balance AS updated_balance
		)
		SELECT
				it.transaction_id,
				it.account_id,
				it.operation_type_id,
				it.amount,
				it.event_date,
				it.balance_before,
				it.balance_after,
				it.created_at,
				it.updated_at
		FROM
				updated_account ua
		JOIN
				inserted_transaction it ON ua.account_id = it.account_id;
	`
	newTxn := models.Transaction{}
	err := r.db.QueryRowContext(ctx, query, txn.AccountID, txn.Amount, txn.OpTypeID).Scan(
		&newTxn.ID,
		&newTxn.AccountID,
		&newTxn.OpTypeID,
		&newTxn.Amount,
		&newTxn.EventDate,
		&newTxn.BalanceBefore,
		&newTxn.BalanceAfter,
		&newTxn.CreatedAt,
		&newTxn.UpdatedAt,
	)

	return newTxn, err
}

func (r *txnRepository) GetTransactionByID(ctx context.Context, txnID int) (*models.Transaction, error) {
	const query = `
		SELECT id, account_id, amount, operation_type_id, event_date, balance_before, balance_after, created_at, updated_at FROM transactions
		WHERE id = $1
	`
	txn := &models.Transaction{}
	err := r.db.QueryRowContext(ctx, query, txn.AccountID, txn.Amount, txn.OpTypeID).Scan(
		&txn.ID,
		&txn.AccountID,
		&txn.Amount,
		&txn.OpTypeID,
		&txn.EventDate,
		&txn.BalanceBefore,
		&txn.BalanceAfter,
		&txn.CreatedAt,
		&txn.UpdatedAt,
	)

	return txn, err
}
