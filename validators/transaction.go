package validators

import (
	apperrors "github.com/sundayezeilo/pismo/app-errors"
	"github.com/sundayezeilo/pismo/dto"
)

func ValidateCreateTransactionReq(params *dto.CreateTxnParams) *apperrors.APIError {
	var errorList []string
	if params.AccountID < 1 {
		errorList = append(errorList, "account_id is required")
	}

	if params.OpTypeID == 0 {
		errorList = append(errorList, "operation_type_id is required")
	}

	if params.Amount <= 0 {
		errorList = append(errorList, "amount be greater than 0")
	}

	if len(errorList) > 0 {
		apiErr := apperrors.ErrBadRequest.WithError(errorList).WithMessage("invalid request parameters")
		return apiErr
	}
	return nil
}
