package services

import (
	"context"
	"fmt"
	"log/slog"
	"math"

	"github.com/sundayezeilo/pismo/constants"
	"github.com/sundayezeilo/pismo/dto"
	apierrors "github.com/sundayezeilo/pismo/errors"
	"github.com/sundayezeilo/pismo/models"
	"github.com/sundayezeilo/pismo/repositories"
)

type TransactionService interface {
	CreateTransaction(context.Context, *dto.CreateTxnParams) (*models.Transaction, error)
	GetTransactionByID(context.Context, int) (*models.Transaction, error)
}

type transactionService struct {
	repo          repositories.TxnRepository
	accService    AccountService
	opTypeService OperationTypeService
}

func NewTransactionService(repo repositories.TxnRepository, accService AccountService, opTypeService OperationTypeService) TransactionService {
	return &transactionService{repo, accService, opTypeService}
}

func (srv *transactionService) CreateTransaction(ctx context.Context, txnParams *dto.CreateTxnParams) (*models.Transaction, error) {
	if err := srv.validateTransaction(ctx, txnParams); err != nil {
		return nil, apierrors.ErrBadRequest.WithMessage(err.Error())
	}
	opType := &models.OperationType{}

	if _, err := srv.validateOpTypes(ctx, txnParams.OpTypeID); err != nil {
		return nil, apierrors.ErrBadRequest.WithMessage("invalid operation type")
	}

	newTxn := &models.Transaction{AccountID: txnParams.AccountID, OpTypeID: txnParams.OpTypeID, Amount: txnParams.Amount * float64(srv.getTxnType(opType.OpType))}
	err := srv.repo.CreateTransaction(ctx, newTxn)

	if err != nil {
		slog.Log(ctx, slog.LevelError, "error creating transaction: "+err.Error())
		return nil, apierrors.ErrInternalServerError.WithMessage("error creating transaction")
	}
	return newTxn, nil
}

func (srv *transactionService) GetTransactionByID(ctx context.Context, txnID int) (*models.Transaction, error) {
	txn, err := srv.repo.GetTransactionByID(ctx, txnID)
	if err != nil {
		return nil, apierrors.ErrNotFound.WithMessage("transaction not found")
	}

	return txn, nil
}

func (srv *transactionService) validateTransaction(ctx context.Context, txnParams *dto.CreateTxnParams) error {
	if _, err := srv.accService.GetAccountByID(ctx, txnParams.AccountID); err != nil {
		return fmt.Errorf("invalid account ID")
	}

	roundedAmount := math.Round(txnParams.Amount*100) / 100

	if roundedAmount != txnParams.Amount {
		return fmt.Errorf("invalid amount: must have at most 2 decimal places")
	}

	return nil
}

func (srv *transactionService) validateOpTypes(ctx context.Context, opTypeID int) (*models.OperationType, error) {
	opType, err := srv.opTypeService.GetOpTypeByID(ctx, opTypeID)
	if err != nil {
		return nil, fmt.Errorf("invalid operation_type_id")
	}

	if !opType.ActiveSupport {
		return nil, fmt.Errorf("operation not currently supported")
	}

	return opType, nil
}

func (srv *transactionService) getTxnType(opDir constants.OperationType) int {
	if opDir == "debit" {
		return -1
	}

	return 1
}
