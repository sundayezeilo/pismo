package tests

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type CreateAccountTestCase struct {
	name       string
	method     string
	url        string
	body       string
	wantStatus int
	wantBody   string
}

type GetAccountTestCase struct {
	name       string
	method     string
	url        string
	pathParam  string
	wantStatus int
	wantBody   string
}

func MakeCreateAccountTestCases(tsUrl string) []*CreateAccountTestCase {
	testCases := []*CreateAccountTestCase{
		{
			name:       "invalid request body",
			method:     "POST",
			url:        tsUrl + "/accounts",
			body:       "",
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"message": "invalid request body", "errors": []}`,
		},
		{
			name:       "missing document_number",
			method:     "POST",
			url:        tsUrl + "/accounts",
			body:       `{}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"message": "document_number is required", "errors": []}`,
		},
		{
			name:       "invalid document_number",
			method:     "POST",
			url:        tsUrl + "/accounts",
			body:       `{"document_number": "EE074746aaaaaxx"}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"message": "invalid document_number", "errors": []}`,
		},
		{
			name:       "valid request",
			method:     "POST",
			url:        tsUrl + "/accounts",
			body:       `{"document_number": "12345678900"}`,
			wantStatus: http.StatusCreated,
			wantBody:   `{"account_id": 1, "document_number": "12345678900"}`,
		},
		{
			name:       "duplicate record",
			method:     "POST",
			url:        tsUrl + "/accounts",
			body:       `{"document_number": "12345678900"}`,
			wantStatus: http.StatusConflict,
			wantBody:   `{"message": "account with 12345678900 already exists", "errors": []}`,
		},
	}
	return testCases
}

func MakeGetAccountTestCases(tsUrl string, accID int, documentNumber string) []*GetAccountTestCase {
	validAccData := map[string]interface{}{
		"account_id":      accID,
		"document_number": documentNumber,
	}

	jsonData, err := json.Marshal(validAccData)
	if err != nil {
		log.Fatal("error occured in MakeCreateAccountTestCases function: ", err)
	}

	testCases := []*GetAccountTestCase{
		{
			name:       "invalid account id",
			method:     "GET",
			url:        tsUrl + "/accounts",
			pathParam:  "wccggAAZZkk",
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"message": "invalid account_id", "errors": []}`,
		},
		{
			name:       "account not found record",
			method:     "GET",
			url:        tsUrl + "/accounts",
			pathParam:  "200",
			wantStatus: http.StatusNotFound,
			wantBody:   `{"message": "account not found", "errors": []}`,
		},
		{
			name:       "valid request",
			method:     "GET",
			url:        tsUrl + "/accounts",
			pathParam:  strconv.Itoa(accID),
			wantStatus: http.StatusOK,
			wantBody:   string(jsonData),
		},
	}
	return testCases
}
