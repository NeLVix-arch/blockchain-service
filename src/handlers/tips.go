package handlers

import (
	"blockchain-service/src/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

// TransferTipHandler is a handler function that transfers a tip from one wallet to another
func TransferTipHandler(DB *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode the request body
		var tip models.Tip
		err := json.NewDecoder(r.Body).Decode(&tip)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Validate the from and to wallet ids
		var fromWallet, toWallet models.Wallet
		err = DB.Get(&fromWallet, "SELECT * FROM wallets WHERE id = $1", tip.FromWalletID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = DB.Get(&toWallet, "SELECT * FROM wallets WHERE id = $1", tip.ToWalletID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Validate the amount
		if tip.Amount <= 0 {
			http.Error(w, "amount must be positive", http.StatusBadRequest)
			return
		}

		// Convert the amount to 18 decimals
		tip.Amount = tip.Amount * 1000000000000000000

		// Start a database transaction
		tx, err := DB.Beginx()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Lock the from and to wallet rows for update
		err = tx.Get(&fromWallet, "SELECT * FROM wallets WHERE id = $1 FOR UPDATE", tip.FromWalletID)
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tx.Get(&toWallet, "SELECT * FROM wallets WHERE id = $1 FOR UPDATE", tip.ToWalletID)
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Decrease the balance of the from wallet
		fromWallet.Balance = fromWallet.Balance - tip.Amount
		// Check for negative balance
		if fromWallet.Balance < 0 {
			tx.Rollback()
			http.Error(w, "insufficient balance", http.StatusBadRequest)
			return
		}

		// Increase the balance of the to wallet
		toWallet.Balance = toWallet.Balance + tip.Amount

		// Update the wallets in the database
		_, err = tx.Exec("UPDATE wallets SET balance = $1 WHERE id = $2", fromWallet.Balance, fromWallet.ID)
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = tx.Exec("UPDATE wallets SET balance = $1 WHERE id = $2", toWallet.Balance, toWallet.ID)
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create the transactions for the tip transfer
		var fromTransaction, toTransaction models.Transaction
		fromTransaction.TransactionType = 2 // BalanceDecrease
		fromTransaction.WalletID = fromWallet.ID
		fromTransaction.Amount = tip.Amount
		fromTransaction.UpdatedBalance = fromWallet.Balance
		toTransaction.TransactionType = 1 // BalanceIncrease
		toTransaction.WalletID = toWallet.ID
		toTransaction.Amount = tip.Amount
		toTransaction.UpdatedBalance = toWallet.Balance

		// Insert the transactions into the database
		err = tx.QueryRowx("INSERT INTO transactions (transaction_type, wallet_id, amount, updated_balance) VALUES ($1, $2, $3, $4) RETURNING id", fromTransaction.TransactionType, fromTransaction.WalletID, fromTransaction.Amount, fromTransaction.UpdatedBalance).Scan(&fromTransaction.ID)
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tx.QueryRowx("INSERT INTO transactions (transaction_type, wallet_id, amount, updated_balance) VALUES ($1, $2, $3, $4) RETURNING id", toTransaction.TransactionType, toTransaction.WalletID, toTransaction.Amount, toTransaction.UpdatedBalance).Scan(&toTransaction.ID)
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

		// Encode the transactions as JSON and send them as response
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode([]models.Transaction{fromTransaction, toTransaction})
		if err != nil {
			log.Println(err)
		}
	}
}
