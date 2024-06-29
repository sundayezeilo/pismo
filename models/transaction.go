package models

import "time"

type Transaction struct {
	ID            int       `json:"transaction_id"`
	AccountID     int       `json:"account_id"`
	OpTypeID      int       `json:"operation_type_id"`
	Amount        float64   `json:"amount"`
	BalanceBefore float64   `json:"balance_before"`
	BalanceAfter  float64   `json:"balance_after"`
	EventDate     time.Time `json:"event_date"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreatedAt     time.Time `json:"created_at"`
}
