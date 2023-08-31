package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB
var TABLES = []string{"Users", "Segments", "UsersSegments"}

func Init(test bool) {
	url := fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable",
		os.Getenv("POSTGRES_DEV_USER"),
		os.Getenv("POSTGRES_DEV_PASSWORD"),
		os.Getenv("POSTGRES_DEV_NAME"),
		os.Getenv("POSTGRES_DEV_DB"))

	if test {
		url = fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable",
			os.Getenv("POSTGRES_TEST_USER"),
			os.Getenv("POSTGRES_TEST_PASSWORD"),
			os.Getenv("POSTGRES_TEST_NAME"),
			os.Getenv("POSTGRES_TEST_DB"))
	}

	db, err := sql.Open(
		"postgres",
		url,
	)
	if err != nil {
		log.Fatalf("sql.Open: %v\n", err)
	}
	DB = db

	if err = DB.Ping(); err != nil {
		log.Fatalf("DB.Ping(): %v\n", err)
	}

	log.Println("You connected to your database.")
}
