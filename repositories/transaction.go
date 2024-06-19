package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/sundayezeilo/pismo/models"

	"github.com/sundayezeilo/pismo/dto"
)

type TxnRepository interface {
	CreateTransaction(context.Context, *dto.CreateTransaction) error
	GetTransactionByID(context.Context, int) (*models.Transaction, error)
}

type txnRepository struct {
	db *sql.DB
}

func NewTxnRepository(db *sql.DB) TxnRepository {
	return &txnRepository{db}
}

func (r *txnRepository) CreateTransaction(ctx context.Context, txn *dto.CreateTransaction) error {
	query := `
		WITH updated_account AS (
				UPDATE accounts
				SET credit_limit = credit_limit + $2
				WHERE id = $1
				RETURNING id AS account_id, credit_limit
		),
		inserted_transaction AS (
				INSERT INTO transactions (account_id, amount, operation_type_id)
				VALUES ($1, $2, $3)
				RETURNING id AS transaction_id, account_id, operation_type_id, amount, event_date, created_at, updated_at
		)
		SELECT
				it.transaction_id,
				it.operation_type_id,
				it.amount,
				it.event_date,
				it.created_at,
				it.updated_at,
				ua.credit_limit
		FROM
				updated_account ua
		JOIN
				inserted_transaction it ON ua.account_id = it.account_id;
	`

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("unable to start transaction: %v", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			slog.Log(ctx, slog.LevelError, "transaction rollback")
			return
		}
	}()

	err = tx.QueryRowContext(ctx, query, txn.AccountID, txn.Amount, txn.OpTypeID).Scan(&txn.TransactionID, &txn.OpTypeID, &txn.Amount, &txn.EventDate, &txn.CreatedAt, &txn.UpdatedAt, &txn.CreditLimit)

	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (r *txnRepository) GetTransactionByID(ctx context.Context, txnID int) (*models.Transaction, error) {
	const query = `
		SELECT id, account_id, amount, operation_type_id, event_date, created_at, updated_at FROM transactions
		WHERE id = $1
	`
	txn := &models.Transaction{}
	err := r.db.QueryRowContext(ctx, query, txnID).
		Scan(&txn.ID, &txn.AccountID, &txn.Amount, &txn.OpTypeID, &txn.EventDate, &txn.CreatedAt, &txn.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return txn, nil
}
