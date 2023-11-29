package handlers

import (
	"blockchain-service/src/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

// CreateWalletHandler is a handler function that creates a new wallet
func CreateWalletHandler(DB *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode the request body
		var wallet models.Wallet
		err := json.NewDecoder(r.Body).Decode(&wallet)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Validate the coin name
		if wallet.Coin == "" {
			http.Error(w, "coin is required", http.StatusBadRequest)
			return
		}

		// fix
		if wallet.Balance != 0 {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		// Insert the wallet into the database
		err = DB.QueryRowx("INSERT INTO wallets (coin, balance) VALUES ($1, $2) RETURNING id", wallet.Coin, 0).Scan(&wallet.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Encode the wallet as JSON and send it as response
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(wallet)
		if err != nil {
			log.Println(err)
		}
	}
}
