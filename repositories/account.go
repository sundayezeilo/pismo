package repositories

import (
	"context"
	"database/sql"

	"github.com/sundayezeilo/pismo/models"
)

type AccountRepository interface {
	CreateAccount(context.Context, *CreateAccountParams) (models.Account, error)
	GetAccountByID(context.Context, int) (models.Account, error)
	GetAccountByDocumentNumber(context.Context, string) (models.Account, error)
}

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &accountRepository{db}
}

type CreateAccountParams struct {
	DocumentNumber string `json:"document_number"`
}

func (r *accountRepository) CreateAccount(ctx context.Context, acc *CreateAccountParams) (models.Account, error) {
	const query = `
		INSERT INTO accounts (document_number)
		VALUES ($1)
		RETURNING id, document_number, balance, created_at, updated_at;
	`
	newAcc := models.Account{}
	err := r.db.QueryRowContext(ctx, query, acc.DocumentNumber).Scan(&newAcc.ID, &newAcc.DocumentNumber, &newAcc.Balance, &newAcc.CreatedAt, &newAcc.UpdatedAt)

	return newAcc, err
}

func (r *accountRepository) GetAccountByID(ctx context.Context, accID int) (models.Account, error) {
	const query = `
		SELECT id, document_number, balance, created_at, updated_at FROM accounts WHERE id = $1;
	`
	acc := models.Account{}
	err := r.db.QueryRowContext(ctx, query, accID).Scan(&acc.ID, &acc.DocumentNumber, &acc.Balance, &acc.CreatedAt, &acc.UpdatedAt)
	return acc, err
}

func (r *accountRepository) GetAccountByDocumentNumber(ctx context.Context, docNum string) (models.Account, error) {
	query := `
		SELECT id, document_number, created_at, updated_at FROM accounts WHERE document_number = $1;
	`
	acc := models.Account{}
	err := r.db.QueryRowContext(ctx, query, docNum).Scan(&acc.ID, &acc.DocumentNumber, &acc.Balance, &acc.CreatedAt, &acc.UpdatedAt)

	return acc, err
}
