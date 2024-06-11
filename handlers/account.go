package handlers

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/sundayezeilo/pismo/dto"
	apierrors "github.com/sundayezeilo/pismo/errors"
	"github.com/sundayezeilo/pismo/services"
	"github.com/sundayezeilo/pismo/validators"
)

type AccountHandler struct {
	Service services.AccountService
}

func NewAccountHandler(service services.AccountService) *AccountHandler {
	return &AccountHandler{service}
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	params := &dto.CreateAccountParams{}

	if err := json.NewDecoder(r.Body).Decode(params); err != nil {
		apierrors.ErrBadRequest.WithMessage("invalid request body").WriteJSON(w)
		return
	}

	defer r.Body.Close()

	if err := validators.ValidateCreateAccountReq(params.DocumentNumber); err != nil {
		apierrors.ErrBadRequest.WithMessage(err.Error()).WriteJSON(w)
		return
	}

	ctx := context.Background()
	acc, err := h.Service.CreateAccount(ctx, params)
	if err != nil {
		if apiErr, ok := err.(*apierrors.APIError); ok {
			apiErr.WriteJSON(w)
		} else {
			slog.Log(ctx, slog.LevelError, "Error creating new account")
			apierrors.ErrInternalServerError.WithMessage("Unexpected error occurred").WriteJSON(w)
		}
		return
	}

	accountResponse := dto.CreateGetAccountResponse{
		AccountID:      acc.ID,
		DocumentNumber: acc.DocumentNumber,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(accountResponse)
}

func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	accIDStr := r.PathValue("accountId")
	log.Println("accIDStr: ", accIDStr)
	accID, err := strconv.Atoi(accIDStr)

	if err != nil || accID < 1 {
		apierrors.ErrBadRequest.WithMessage("invalid account_id").WriteJSON(w)
		return
	}

	ctx := context.Background()
	account, err := h.Service.GetAccountByID(ctx, accID)

	if err != nil {
		apierrors.ErrNotFound.WithMessage("account not found").WriteJSON(w)
		return
	}

	accountResponse := dto.CreateGetAccountResponse{
		AccountID:      account.ID,
		DocumentNumber: account.DocumentNumber,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accountResponse)
}
