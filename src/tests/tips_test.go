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

// TestTransferTipHandler is a test function for the transfer tip handler
func TestTransferTipHandler(t *testing.T) {
	db := InitTESTdb()
	defer db.Close()

	// Prepare the request body
	tip := models.Tip{FromWalletID: 1, ToWalletID: 2, Amount: 0.123}
	body, err := json.Marshal(tip)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new request
	req, err := http.NewRequest("POST", "/transfer_tips", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler function
	handlers.TransferTipHandler(db)(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body
	var response []models.Transaction
	json.Unmarshal(rr.Body.Bytes(), &response)

	assert.Equal(t, tip.FromWalletID, int64(response[1].TransactionType))
	assert.Equal(t, tip.ToWalletID, int64(response[0].TransactionType))
	assert.Equal(t, tip.Amount*1000000000000000000, response[0].Amount)
	assert.Equal(t, tip.Amount*1000000000000000000, response[1].Amount)
	assert.Equal(t, tip.FromWalletID, response[0].WalletID)
	assert.Equal(t, tip.ToWalletID, response[1].WalletID)
}
