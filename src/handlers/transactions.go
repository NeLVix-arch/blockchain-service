package handlers

import (
	"blockchain-service/src/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

// ProcessTransactionHandler is a handler function that processes a transaction
func ProcessTransactionHandler(DB *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode the request body
		var transaction models.Transaction
		err := json.NewDecoder(r.Body).Decode(&transaction)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Validate the transaction type
		if transaction.TransactionType != 1 && transaction.TransactionType != 2 {
			http.Error(w, "invalid transaction type", http.StatusBadRequest)
			return
		}

		// Validate the wallet id
		var wallet models.Wallet
		err = DB.Get(&wallet, "SELECT * FROM wallets WHERE id = $1", transaction.WalletID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Validate the amount
		if transaction.Amount <= 0 {
			http.Error(w, "amount must be positive", http.StatusBadRequest)
			return
		}

		// Convert the amount to 18 decimals
		transaction.Amount = transaction.Amount * 1000000000000000000

		// Start a database transaction
		tx, err := DB.Beginx()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Lock the wallet row for update
		err = tx.Get(&wallet, "SELECT * FROM wallets WHERE id = $1 FOR UPDATE", transaction.WalletID)
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Update the wallet balance
		if transaction.TransactionType == 1 {
			// Increase the balance
			wallet.Balance = wallet.Balance + transaction.Amount
		} else {
			// Decrease the balance
			wallet.Balance = wallet.Balance - transaction.Amount
			// Check for negative balance
			if wallet.Balance < 0 {
				tx.Rollback()
				http.Error(w, "insufficient balance", http.StatusBadRequest)
				return
			}
		}

		// Update the wallet in the database
		_, err = tx.Exec("UPDATE wallets SET balance = $1 WHERE id = $2", wallet.Balance, wallet.ID)
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the updated balance in the transaction
		transaction.UpdatedBalance = wallet.Balance

		// Insert the transaction into the database
		err = tx.QueryRowx("INSERT INTO transactions (transaction_type, wallet_id, amount, updated_balance) VALUES ($1, $2, $3, $4) RETURNING id", transaction.TransactionType, transaction.WalletID, transaction.Amount, transaction.UpdatedBalance).Scan(&transaction.ID)
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Commit the database transaction
		err = tx.Commit()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Encode the transaction as JSON and send it as response
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(transaction)
		if err != nil {
			log.Println(err)
		}
	}
}
