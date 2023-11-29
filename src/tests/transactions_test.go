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

// TestProcessTransactionHandler is a test function for the process transaction handler
func TestProcessTransactionHandler(t *testing.T) {
	db := InitTESTdb()
	defer db.Close()

	// Prepare the request body
	transaction := models.Transaction{TransactionType: 1, WalletID: 1, Amount: 0.123}
	body, err := json.Marshal(transaction)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new request
	req, err := http.NewRequest("POST", "/process_transaction", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler function
	handlers.ProcessTransactionHandler(db)(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body
	var processedTransaction models.Transaction
	json.Unmarshal(rr.Body.Bytes(), &processedTransaction)
	assert.NotZero(t, processedTransaction.ID)
	assert.Equal(t, transaction.TransactionType, processedTransaction.TransactionType)
	assert.Equal(t, transaction.WalletID, processedTransaction.WalletID)
	assert.Equal(t, transaction.Amount*1000000000000000000, processedTransaction.Amount)
}
