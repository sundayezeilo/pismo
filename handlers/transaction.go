package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/sundayezeilo/pismo/dto"
	apierrors "github.com/sundayezeilo/pismo/errors"
	"github.com/sundayezeilo/pismo/services"
	"github.com/sundayezeilo/pismo/validators"
)

type TransactionHandler struct {
	Service services.TransactionService
}

func NewTransactionHandler(service services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service}
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	params := &dto.CreateTxnParams{}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		apierrors.ErrBadRequest.WithMessage("invalid request payload").WriteJSON(w)
		return
	}

	defer r.Body.Close()

	if err := validators.ValidateCreateTransactionReq(params); err != nil {
		err.WriteJSON(w)
		return
	}

	ctx := context.Background()
	txn, err := h.Service.CreateTransaction(ctx, params)

	if err != nil {
		if apiErr, ok := err.(*apierrors.APIError); ok {
			apiErr.WriteJSON(w)
		} else {
			slog.Log(ctx, slog.LevelError, "Error creating new transaction")
			apierrors.ErrInternalServerError.WithMessage("Unexpected error occurred").WriteJSON(w)
		}
		return
	}

	txnResponse := dto.CreateTxnResponse{
		TransactionID: txn.ID,
		AccountID:     txn.AccountID,
		OpTypeID:      txn.OpTypeID,
		Amount:        txn.Amount,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(txnResponse)
}
