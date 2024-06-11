package repositories

import (
	"context"
	"database/sql"

	"github.com/sundayezeilo/pismo/models"
)

type AccountRepository interface {
	CreateAccount(context.Context, *models.Account) error
	GetAccountByID(context.Context, int) (*models.Account, error)
	GetAccountByDocumentNumber(context.Context, string) (*models.Account, error)
}

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &accountRepository{db}
}

func (r *accountRepository) CreateAccount(ctx context.Context, acc *models.Account) error {
	const query = `
		INSERT INTO accounts (document_number)
		VALUES ($1)
		RETURNING *;
	`
	err := r.db.QueryRowContext(ctx, query, acc.DocumentNumber).Scan(&acc.ID, &acc.DocumentNumber, &acc.CreatedAt, &acc.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (r *accountRepository) GetAccountByID(ctx context.Context, accID int) (*models.Account, error) {
	const query = `
		SELECT * FROM accounts WHERE id = $1;
	`
	acc := &models.Account{}
	err := r.db.QueryRowContext(ctx, query, accID).Scan(&acc.ID, &acc.DocumentNumber, &acc.CreatedAt, &acc.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

func (r *accountRepository) GetAccountByDocumentNumber(ctx context.Context, docNum string) (*models.Account, error) {
	query := `
		SELECT * FROM accounts WHERE document_number = $1;
	`
	acc := &models.Account{}
	err := r.db.QueryRowContext(ctx, query, docNum).Scan(&acc.ID, &acc.DocumentNumber, &acc.CreatedAt, &acc.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return acc, nil
}
