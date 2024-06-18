package apierrors

import (
	"encoding/json"
	"net/http"

	"github.com/sundayezeilo/pismo/dto"
)

type APIError struct {
	Code    int      `json:"-"`
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

func NewAPIError(code int, message string, errorList []string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Errors:  errorList,
	}
}

func (e *APIError) WithMessage(message string) *APIError {
	return &APIError{
		Code:    e.Code,
		Message: message,
		Errors:  e.Errors,
	}
}

func (e *APIError) WithError(errorList []string) *APIError {
	return &APIError{
		Code:    e.Code,
		Message: e.Message,
		Errors:  errorList,
	}
}

func (e *APIError) Error() string {
	return e.Message
}

func (e *APIError) WriteJSON(w http.ResponseWriter) {
	resp := dto.ErrorResponse{
		Status:  false,
		Message: e.Message,
		Errors:  e.Errors,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Code)
	json.NewEncoder(w).Encode(resp)
}

var (
	ErrBadRequest          = NewAPIError(http.StatusBadRequest, "Bad request", []string{})
	ErrConflict            = NewAPIError(http.StatusConflict, "Resource already exists", []string{})
	ErrNotFound            = NewAPIError(http.StatusNotFound, "Resource not found", []string{})
	ErrInternalServerError = NewAPIError(http.StatusInternalServerError, "Internal server error", []string{})
	ErrUnauthorized        = NewAPIError(http.StatusUnauthorized, "Unauthorized", []string{})
	ErrForbidden           = NewAPIError(http.StatusForbidden, "Forbidden", []string{})
)
