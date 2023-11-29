package database

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DB is a global variable that holds the database connection
var DB *sqlx.DB

// GORM is a global variable that holds the gorm ORM instance
var GORM *gorm.DB

var osenv = os.Getenv("POSTGRES_CONFIG")

func InitDBS() (*sqlx.DB, *gorm.DB) {
	// Connect to the database
	var err error
	DB, err = sqlx.Connect("postgres", osenv)
	if err != nil {
		log.Fatal(err)
	}

	// make
	DB.Exec(`CREATE TABLE IF NOT EXISTS wallets (
		id SERIAL PRIMARY KEY,
		coin VARCHAR(20) NOT NULL,
		balance NUMERIC(38, 0) NOT NULL
	);`)

	// Initialize gorm
	GORM, err = gorm.Open("postgres", osenv)
	if err != nil {
		DB.Close()
		log.Fatal(err)
	}

	return DB, GORM
}
