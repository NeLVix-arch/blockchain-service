package tests

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var osenvtest = os.Getenv("POSTGRES_CONFIG_TESTING")

func InitTESTdb() *sqlx.DB {
	// Connect to the database
	DB, err := sqlx.Connect("postgres", osenvtest)
	if err != nil {
		log.Fatal(err)
	}
	return DB
}
