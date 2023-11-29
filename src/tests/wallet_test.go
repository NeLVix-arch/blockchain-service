package tests

import (
	"blockchain-service/src/handlers"
	"blockchain-service/src/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateWalletHandler(t *testing.T) {
	DB := InitTESTdb()
	// Prepare the test cases
	testCases := []struct {
		name       string
		input      string
		statusCode int
		output     models.Wallet
	}{
		{
			name:       "valid input",
			input:      "{\"coin\": \"eth\"}",
			statusCode: http.StatusOK,
			output:     models.Wallet{ID: 1, Coin: "eth", Balance: 5},
		},
		{
			name:       "empty input",
			input:      "{}",
			statusCode: http.StatusBadRequest,
			output:     models.Wallet{},
		},
		{
			name:       "invalid input",
			input:      "{\"coin\": 123}",
			statusCode: http.StatusBadRequest,
			output:     models.Wallet{},
		},
		{
			name:       "empty coin",
			input:      "{\"coin\": \"\"}",
			statusCode: http.StatusBadRequest,
			output:     models.Wallet{},
		},
	}

	// Iterate over the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new request with the input
			req, err := http.NewRequest("POST", "/create_wallet", bytes.NewBufferString(tc.input))
			if err != nil {
				t.Fatal(err)
			}

			// Create a new response recorder
			rr := httptest.NewRecorder()

			// Call the handler function
			handlers.CreateWalletHandler(DB)(rr, req)

			// Check the status code
			assert.Equal(t, tc.statusCode, rr.Code)

			// Check the output
			var wallet models.Wallet
			json.NewDecoder(rr.Body).Decode(&wallet)
			assert.Equal(t, float64(0), wallet.Balance)
		})
	}
}
