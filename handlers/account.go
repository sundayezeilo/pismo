package handlers

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	apperrors "github.com/sundayezeilo/pismo/app-errors"
	"github.com/sundayezeilo/pismo/dto"
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
	params := &dto.CreateAccountRequest{}

	if err := json.NewDecoder(r.Body).Decode(params); err != nil {
		apperrors.NewAPIError(http.StatusBadRequest, "invalid request body").WriteJSON(w)
		return
	}

	defer r.Body.Close()

	if err := validators.ValidateCreateAccountReq(params.DocumentNumber); err != nil {
		apperrors.NewAPIError(http.StatusBadRequest, err.Error()).WriteJSON(w)
		return
	}

	acc, err := h.Service.CreateAccount(r.Context(), params)
	if err != nil {
		if apiErr, ok := err.(*apperrors.APIError); ok {
			apiErr.WriteJSON(w)
		} else {
			slog.Log(r.Context(), slog.LevelError, "Error creating new account")
			apperrors.NewAPIError(http.StatusInternalServerError, "Unexpected error occurred").WriteJSON(w)
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
		apperrors.NewAPIError(http.StatusBadRequest, "invalid account_id").WriteJSON(w)
		return
	}

	account, err := h.Service.GetAccountByID(r.Context(), accID)

	if err != nil {
		apperrors.NewAPIError(http.StatusBadRequest, "account not found").WriteJSON(w)
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
