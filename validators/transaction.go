package validators

import (
	"github.com/sundayezeilo/pismo/dto"
	apierrors "github.com/sundayezeilo/pismo/errors"
)

func ValidateCreateTransactionReq(params *dto.CreateTxnParams) *apierrors.APIError {
	var errorList []string
	if params.AccountID < 1 {
		errorList = append(errorList, "account_id is required")
	}

	if params.OpTypeID == 0 {
		errorList = append(errorList, "operation_type_id is required")
	}

	if params.Amount == 0 {
		errorList = append(errorList, "amount can not be 0")
	}
	if len(errorList) > 0 {
		apiErr := apierrors.ErrBadRequest.WithError(errorList).WithMessage("invalid request parameters")
		return apiErr
	}
	return nil
}
