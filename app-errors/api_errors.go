package apperrors

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

func NewAPIError(code int, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Errors:  []string{},
	}
}

func (e *APIError) WithMessage(message string) *APIError {
	e.Message = message
	return e
}

func (e *APIError) WithError(errorList []string) *APIError {
	e.Errors = errorList
	return e
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
