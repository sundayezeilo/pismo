package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sundayezeilo/pismo/dto"
	apierrors "github.com/sundayezeilo/pismo/errors"
	"github.com/sundayezeilo/pismo/models"
	"github.com/sundayezeilo/pismo/repositories"
)

type AccountService interface {
	CreateAccount(context.Context, *dto.CreateAccountParams) (*models.Account, error)
	GetAccountByID(context.Context, int) (*models.Account, error)
}

type accountService struct {
	repo repositories.AccountRepository
}

func NewAccountService(repo repositories.AccountRepository) AccountService {
	return &accountService{repo}
}

func (srv *accountService) CreateAccount(ctx context.Context, params *dto.CreateAccountParams) (*models.Account, error) {
	_, err := srv.repo.GetAccountByDocumentNumber(ctx, params.DocumentNumber)
	if err == nil {
		apiErr := apierrors.ErrConflict.WithMessage(fmt.Sprintf("account with %v already exists", params.DocumentNumber))
		return nil, apiErr
	}
	newAcc := &models.Account{DocumentNumber: params.DocumentNumber}
	err = srv.repo.CreateAccount(ctx, newAcc)

	if err != nil {
		slog.Log(ctx, slog.LevelError, "error creating new account: "+err.Error())
		return nil, apierrors.ErrInternalServerError.WithMessage("error creating new account")
	}
	return newAcc, nil
}

func (srv *accountService) GetAccountByID(ctx context.Context, accID int) (*models.Account, error) {
	acc, err := srv.repo.GetAccountByID(ctx, accID)

	if err != nil {
		apiErr := apierrors.ErrNotFound.WithMessage(fmt.Sprintf("no account found with ID: %v", accID))
		return nil, apiErr
	}
	return acc, nil
}
