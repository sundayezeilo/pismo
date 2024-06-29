package dto

type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number"`
}

type CreateGetAccountResponse struct {
	AccountID      int     `json:"account_id"`
	DocumentNumber string  `json:"document_number"`
	Balance    float64 `json:"balance"`
}
