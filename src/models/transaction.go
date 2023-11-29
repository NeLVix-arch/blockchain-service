package models

// Transaction is a struct that represents a transaction in the database
type Transaction struct {
	ID              int64   `json:"id" db:"id" gorm:"primary_key"`
	TransactionType int     `json:"transaction_type" db:"transaction_type" gorm:"not null"`
	WalletID        int64   `json:"wallet_id" db:"wallet_id" gorm:"not null"`
	Amount          float64 `json:"amount" db:"amount" gorm:"not null"`
	UpdatedBalance  float64 `json:"updated_balance" db:"updated_balance" gorm:"not null"`
}
