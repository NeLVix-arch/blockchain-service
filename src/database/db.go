package database

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DB is a global variable that holds the database connection
var DB *sqlx.DB

// GORM is a global variable that holds the gorm ORM instance
var GORM *gorm.DB

func InitDBS() (*sqlx.DB, *gorm.DB) {
	// Connect to the database
	var err error
	DB, err = sqlx.Connect("postgres", "user=postgres password=1111 dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize gorm
	GORM, err = gorm.Open("postgres", "user=postgres password=1111 dbname=postgres sslmode=disable")
	if err != nil {
		DB.Close()
		log.Fatal(err)
	}

	return DB, GORM
}
