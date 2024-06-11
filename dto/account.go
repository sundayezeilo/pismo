package dto

type CreateAccountParams struct {
	DocumentNumber string `json:"document_number"`
}

type CreateGetAccountResponse struct {
	AccountID      int    `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}
