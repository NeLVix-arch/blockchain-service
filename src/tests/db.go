package tests

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func InitTESTdb() *sqlx.DB {
	// Connect to the database
	DB, err := sqlx.Connect("postgres", "user=postgres password=1111 dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return DB
}
