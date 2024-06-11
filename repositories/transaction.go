package repositories

import (
	"context"
	"database/sql"

	"github.com/sundayezeilo/pismo/models"
)

type TxnRepository interface {
	CreateTransaction(context.Context, *models.Transaction) error
	GetTransactionByID(context.Context, int) (*models.Transaction, error)
}

type txnRepository struct {
	db *sql.DB
}

func NewTxnRepository(db *sql.DB) TxnRepository {
	return &txnRepository{db}
}

func (r *txnRepository) CreateTransaction(ctx context.Context, txn *models.Transaction) error {

	const query = `
		INSERT INTO transactions (account_id, operation_type_id, amount)
		VALUES ($1, $2, $3)
		RETURNING *;
	`
	err := r.db.QueryRowContext(ctx, query, txn.AccountID, txn.OpTypeID, txn.Amount).Scan(&txn.ID, &txn.AccountID, &txn.OpTypeID, &txn.Amount, &txn.EventDate, &txn.CreatedAt, &txn.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *txnRepository) GetTransactionByID(ctx context.Context, txnID int) (*models.Transaction, error) {
	const query = `
		SELECT * FROM transactions
		WHERE id = $1
	`
	txn := &models.Transaction{}
	err := r.db.QueryRowContext(ctx, query, txnID).
		Scan(&txn.ID, &txn.AccountID, &txn.OpTypeID, &txn.Amount, &txn.EventDate, &txn.CreatedAt, &txn.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return txn, nil
}
