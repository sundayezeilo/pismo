package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

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

type CreateTxnResponse struct {
	TransactionID int       `json:"transaction_id"`
	AccountID     int       `json:"account_id"`
	OpTypeID      int       `json:"operation_type_id"`
	Amount        float64   `json:"amount"`
	EventDate     time.Time `json:"event_date"`
	BalanceBefore float64   `json:"balance_before"`
	BalanceAfter  float64   `json:"balance_after"`
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	params := &dto.CreateTxnRequest{}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		apperrors.NewAPIError(http.StatusBadRequest, "invalid request payload").WriteJSON(w)
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
			apperrors.NewAPIError(http.StatusInternalServerError, "Unexpected error occurred").WriteJSON(w)
		}
		return
	}

	newTx := CreateTxnResponse{
		TransactionID: txn.TransactionID,
		AccountID:     txn.AccountID,
		OpTypeID:      txn.OpTypeID,
		Amount:        txn.Amount,
		EventDate:     txn.EventDate,
		BalanceBefore: txn.BalanceBefore,
		BalanceAfter:  txn.BalanceAfter,
	}

	resp := dto.SuccessResponse{
		Status:  true,
		Message: "transaction successfully created",
		Data:    newTx,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
