package models

import (
	"time"

	"github.com/sundayezeilo/pismo/constants"
)

type OperationType struct {
	ID            int                                `json:"id"`
	Description   constants.OperationTypeDescription `json:"description"`
	OpType        constants.OperationType            `json:"op_type"`
	ActiveSupport bool                               `json:"active_support"`
	CreatedAt     time.Time                          `json:"created_at"`
	UpdatedAt     time.Time                          `json:"updated_at"`
}
