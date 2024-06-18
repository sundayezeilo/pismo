package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	apperrors "github.com/sundayezeilo/pismo/app-errors"
	"github.com/sundayezeilo/pismo/dto"
	"github.com/sundayezeilo/pismo/services"
	"github.com/sundayezeilo/pismo/validators"
)

type TransactionHandler struct {
	service services.TransactionService
}

func NewTransactionHandler(service services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service}
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	params := &dto.CreateTxnParams{}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		apperrors.ErrBadRequest.WithMessage("invalid request payload").WriteJSON(w)
		return
	}

	defer r.Body.Close()

	if err := validators.ValidateCreateTransactionReq(params); err != nil {
		err.WriteJSON(w)
		return
	}

	txn, err := h.service.CreateTransaction(r.Context(), params)

	if err != nil {
		if apiErr, ok := err.(*apperrors.APIError); ok {
			apiErr.WriteJSON(w)
		} else {
			slog.Log(r.Context(), slog.LevelError, "Error creating new transaction")
			apperrors.ErrInternalServerError.WithMessage("Unexpected error occurred").WriteJSON(w)
		}
		return
	}

	newTx := dto.CreateTxnResponse{
		TransactionID: txn.TransactionID,
		AccountID:     txn.AccountID,
		OpTypeID:      txn.OpTypeID,
		Amount:        txn.Amount,
		EventDate:     txn.EventDate,
		CreditLimit:   txn.CreditLimit,
	}

	resp := dto.SuccessResponse{
		Status:  true,
		Message: "transaction successfully created",
		Data:    newTx,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
