package dto

type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number"`
}

type CreateGetAccountResponse struct {
	AccountID      int     `json:"account_id"`
	DocumentNumber string  `json:"document_number"`
	CreditLimit    float64 `json:"credit_limit"`
}
