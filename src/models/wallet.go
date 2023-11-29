package models

// Wallet is a struct that represents a wallet in the database
type Wallet struct {
	ID      int64   `json:"id" db:"id" gorm:"primary_key"`
	Coin    string  `json:"coin" db:"coin" gorm:"not null"`
	Balance float64 `json:"balance" db:"balance" gorm:"not null"`
}
