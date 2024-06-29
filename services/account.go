package services

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	apperrors "github.com/sundayezeilo/pismo/app-errors"
	"github.com/sundayezeilo/pismo/dto"
	"github.com/sundayezeilo/pismo/repositories"
)

type AccountService interface {
	CreateAccount(context.Context, *dto.CreateAccountRequest) (accountResponse, error)
	GetAccountByID(context.Context, int) (accountResponse, error)
}

type accountService struct {
	repo repositories.AccountRepository
}

func NewAccountService(repo repositories.AccountRepository) AccountService {
	return &accountService{repo}
}

type accountResponse struct {
	ID             int       `json:"id"`
	DocumentNumber string    `json:"document_number"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Balance        float64   `json:"balance"`
}

func (srv *accountService) CreateAccount(ctx context.Context, params *dto.CreateAccountRequest) (accountResponse, error) {
	resp := accountResponse{}
	_, err := srv.repo.GetAccountByDocumentNumber(ctx, params.DocumentNumber)

	if err == nil {
		apiErr := apperrors.NewAPIError(http.StatusConflict, fmt.Sprintf("account with %v already exists", params.DocumentNumber))
		return resp, apiErr
	}

	acc := &repositories.CreateAccountParams{DocumentNumber: params.DocumentNumber}
	newAcc, err := srv.repo.CreateAccount(ctx, acc)

	if err != nil {
		slog.Log(ctx, slog.LevelError, "error creating new account: "+err.Error())
		return resp, apperrors.NewAPIError(http.StatusInternalServerError, "error creating new account")
	}

	resp = accountResponse{
		ID:             newAcc.ID,
		DocumentNumber: newAcc.DocumentNumber,
		CreatedAt:      newAcc.CreatedAt,
		UpdatedAt:      newAcc.UpdatedAt,
		Balance:        newAcc.Balance,
	}

	return resp, nil
}

func (srv *accountService) GetAccountByID(ctx context.Context, accID int) (accountResponse, error) {
	resp := accountResponse{}
	acc, err := srv.repo.GetAccountByID(ctx, accID)

	if err != nil {
		apiErr := apperrors.NewAPIError(http.StatusNotFound, fmt.Sprintf("no account found with ID: %v", accID))
		return resp, apiErr
	}
	resp = accountResponse{
		ID:             acc.ID,
		DocumentNumber: acc.DocumentNumber,
		CreatedAt:      acc.CreatedAt,
		UpdatedAt:      acc.UpdatedAt,
		Balance:        acc.Balance,
	}

	return resp, nil
}
