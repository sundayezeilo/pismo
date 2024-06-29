package dto

type CreateTxnRequest struct {
	AccountID int     `json:"account_id" validate:"required,number"`
	OpTypeID  int     `json:"operation_type_id" validate:"required,number"`
	Amount    float64 `json:"amount" validate:"required,number"`
}
