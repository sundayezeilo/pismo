package tests

import (
	"encoding/json"
	"log"
	"net/http"
)

type CreateTransactionTestCase struct {
	name       string
	method     string
	url        string
	body       string
	wantStatus int
}

func MakeCreateTransactionTestCases(
	tsUrl string,
	accID int,
	opTypeID int,
	amount float64,
) []*CreateTransactionTestCase {

	correctPayload := map[string]interface{}{
		"account_id":        accID,
		"operation_type_id": opTypeID,
		"amount":            amount,
	}

	correctAmountJsonData, err := json.Marshal(correctPayload)
	if err != nil {
		log.Fatal("error occured in MakeCreateTransactionTestCases function: ", err)
	}

	wrongAmountPayload := map[string]interface{}{
		"account_id":        accID,
		"operation_type_id": opTypeID,
		"amount":            25.7347,
	}

	wrongAmountJsonData, err := json.Marshal(wrongAmountPayload)
	if err != nil {
		log.Fatal("error occured in MakeCreateTransactionTestCases function: ", err)
	}

	testCases := []*CreateTransactionTestCase{
		{
			name:       "missing request body",
			method:     "POST",
			url:        tsUrl + "/transactions",
			body:       "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid request payload",
			method:     "POST",
			url:        tsUrl + "/transactions",
			body:       `{}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid params",
			method:     "POST",
			url:        tsUrl + "/transactions",
			body:       `{"account_id":"invalid","operation_type_id":"2","amount":0}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid amount in the request",
			method:     "POST",
			url:        tsUrl + "/transactions",
			body:       string(wrongAmountJsonData),
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "valid request in the request",
			method:     "POST",
			url:        tsUrl + "/transactions",
			body:       string(correctAmountJsonData),
			wantStatus: http.StatusCreated,
		},
	}
	return testCases
}
