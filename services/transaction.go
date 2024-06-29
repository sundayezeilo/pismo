package services

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"time"

	apperrors "github.com/sundayezeilo/pismo/app-errors"
	"github.com/sundayezeilo/pismo/constants"
	"github.com/sundayezeilo/pismo/dto"
	"github.com/sundayezeilo/pismo/models"
	"github.com/sundayezeilo/pismo/repositories"
)

type TransactionService interface {
	CreateTransaction(context.Context, *dto.CreateTxnRequest) (transactionResponse, error)
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

type transactionResponse struct {
	TransactionID int       `json:"transaction_id"`
	AccountID     int       `json:"account_id"`
	OpTypeID      int       `json:"operation_type_id"`
	Amount        float64   `json:"amount"`
	EventDate     time.Time `json:"event_date"`
	BalanceBefore float64   `json:"balance_before"`
	BalanceAfter  float64   `json:"balance_after"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (srv *transactionService) CreateTransaction(ctx context.Context, txnParams *dto.CreateTxnRequest) (transactionResponse, error) {
	resp := transactionResponse{}
	if err := srv.validateTransaction(ctx, txnParams); err != nil {
		return resp, apperrors.NewAPIError(http.StatusBadRequest, err.Error())
	}
	var opType *models.OperationType

	opType, err := srv.validateOpTypes(ctx, txnParams.OpTypeID)

	if err != nil {
		return resp, apperrors.NewAPIError(http.StatusBadRequest, "invalid operation type")
	}

	usrAcc, err := srv.accService.GetAccountByID(ctx, txnParams.AccountID)
	if err != nil {
		slog.Log(ctx, slog.LevelError, "error creating transaction: "+err.Error())
		return resp, apperrors.NewAPIError(http.StatusInternalServerError, "error creating transaction")
	}

	if usrAcc.Balance < txnParams.Amount && opType.OpType != constants.Credit {
		slog.Log(ctx, slog.LevelError, "insufficient credit")
		return resp, apperrors.NewAPIError(http.StatusBadRequest, "insufficient credit")
	}

	payload := &repositories.CreateTransactionParams{AccountID: txnParams.AccountID, OpTypeID: txnParams.OpTypeID, Amount: txnParams.Amount * float64(srv.getTxnType(opType.OpType))}
	newTxn, err := srv.repo.CreateTransaction(ctx, payload)

	if err != nil {
		slog.Log(ctx, slog.LevelError, "error creating transaction: "+err.Error())
		return resp, apperrors.NewAPIError(http.StatusInternalServerError, "error creating transaction")
	}

	resp = transactionResponse{
		TransactionID: newTxn.ID,
		AccountID:     newTxn.AccountID,
		OpTypeID:      newTxn.OpTypeID,
		Amount:        newTxn.Amount,
		EventDate:     newTxn.EventDate,
		CreatedAt:     newTxn.CreatedAt,
		UpdatedAt:     newTxn.UpdatedAt,
		BalanceBefore: newTxn.BalanceBefore,
		BalanceAfter:  newTxn.BalanceAfter,
	}

	return resp, nil
}

func (srv *transactionService) GetTransactionByID(ctx context.Context, txnID int) (*models.Transaction, error) {
	txn, err := srv.repo.GetTransactionByID(ctx, txnID)
	if err != nil {
		return nil, apperrors.NewAPIError(http.StatusNotFound, "transaction not found")
	}

	return txn, nil
}

func (srv *transactionService) validateTransaction(ctx context.Context, txnParams *dto.CreateTxnRequest) error {
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
