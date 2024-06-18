package handlers

import (
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

	acc, err := h.Service.CreateAccount(r.Context(), params)
	if err != nil {
		if apiErr, ok := err.(*apierrors.APIError); ok {
			apiErr.WriteJSON(w)
		} else {
			slog.Log(r.Context(), slog.LevelError, "Error creating new account")
			apierrors.ErrInternalServerError.WithMessage("Unexpected error occurred").WriteJSON(w)
		}
		return
	}

	newAcc := dto.CreateGetAccountResponse{
		AccountID:      acc.ID,
		DocumentNumber: acc.DocumentNumber,
	}

	resp := dto.SuccessResponse{
		Status:  true,
		Message: "account successfully created",
		Data:    newAcc,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	accIDStr := r.PathValue("accountId")
	log.Println("accIDStr: ", accIDStr)
	accID, err := strconv.Atoi(accIDStr)

	if err != nil || accID < 1 {
		apierrors.ErrBadRequest.WithMessage("invalid account_id").WriteJSON(w)
		return
	}

	account, err := h.Service.GetAccountByID(r.Context(), accID)

	if err != nil {
		apierrors.ErrNotFound.WithMessage("account not found").WriteJSON(w)
		return
	}

	acc := dto.CreateGetAccountResponse{
		AccountID:      account.ID,
		DocumentNumber: account.DocumentNumber,
		CreditLimit:    account.CreditLimit,
	}

	resp := dto.SuccessResponse{
		Status:  true,
		Message: "account successfully retrieved",
		Data:    acc,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
