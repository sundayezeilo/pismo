package models

import "time"

type Transaction struct {
	ID        int       `json:"transaction_id"`
	AccountID int       `json:"account_id"`
	OpTypeID  int       `json:"operation_type_id"`
	Amount    float64   `json:"amount"`
	EventDate time.Time `json:"event_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
