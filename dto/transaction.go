package dto

import "time"

type CreateTxnParams struct {
	AccountID int     `json:"account_id" validate:"required,number"`
	OpTypeID  int     `json:"operation_type_id" validate:"required,number"`
	Amount    float64 `json:"amount" validate:"required,number"`
}

type CreateTxnResponse struct {
	TransactionID int       `json:"transaction_id"`
	AccountID     int       `json:"account_id"`
	OpTypeID      int       `json:"operation_type_id"`
	Amount        float64   `json:"amount"`
	EventDate     time.Time `json:"event_date"`
	CreditLimit   float64   `json:"credit_limit"`
}

type CreateTransaction struct {
	TransactionID int       `json:"transaction_id"`
	AccountID     int       `json:"account_id"`
	OpTypeID      int       `json:"operation_type_id"`
	Amount        float64   `json:"amount"`
	EventDate     time.Time `json:"event_date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreditLimit   float64   `json:"credit_limit"`
}
