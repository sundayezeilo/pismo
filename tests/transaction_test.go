package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sundayezeilo/pismo/app"
	"github.com/sundayezeilo/pismo/config"
)

func TestCreateTransaction(t *testing.T) {
	envFilePath := GetEnvPath()

	cfg := config.LoadEnv(envFilePath)
	app := app.NewApp(cfg)
	defer app.Cleanup()

	ts := httptest.NewServer(app.Router)
	defer ts.Close()

	documentNumber := GenerateRandomNumberString(11)
	accountData := map[string]interface{}{
		"document_number": documentNumber,
	}

	jsonData, _ := json.Marshal(accountData)

	resp, _ := http.Post(ts.URL+"/accounts", "application/json", bytes.NewBuffer(jsonData))
	var respData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respData)
	defer resp.Body.Close()

	accountIDFloat, _ := respData["account_id"].(float64)
	accID := int(accountIDFloat)

	amount := 800.24
	opeType := 2
	testCases := MakeCreateTransactionTestCases(ts.URL, accID, opeType, amount)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.url, bytes.NewBuffer([]byte(tc.body)))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("could not send request: %v", err)
			}
			defer resp.Body.Close()

			assert.Equal(t, tc.wantStatus, resp.StatusCode)
		})
	}
}
