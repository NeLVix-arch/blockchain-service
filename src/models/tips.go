package models

// Tip is a struct that represents a tip transfer request
type Tip struct {
	FromWalletID int64   `json:"from_wallet_id"`
	ToWalletID   int64   `json:"to_wallet_id"`
	Amount       float64 `json:"amount"`
}
