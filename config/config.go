package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

var DB *sql.DB
var TABLES = []string{"Records", "UsersSegments", "Users", "Segments"}

var POSTGRES_USER = os.Getenv("POSTGRES_DEV_USER")
var POSTGRES_PASSWORD = os.Getenv("POSTGRES_DEV_PASSWORD")
var POSTGRES_NAME = os.Getenv("POSTGRES_DEV_NAME")
var POSTGRES_DB = os.Getenv("POSTGRES_DEV_DB")
var REPORTS_DIRNAME = os.Getenv("REPORTS_DEV_DIRNAME")
var GENERATED_DIRNAME = os.Getenv("GENERATED_DIRNAME")

func ensureDir(path string) {
	err := os.MkdirAll(path, os.ModeDir)
	if err != nil {
		log.Fatalf("os.MkdirAll(): %v\n", err)
	}
}

func connectToDB() {
	url := fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable",
		POSTGRES_USER,
		POSTGRES_PASSWORD,
		POSTGRES_NAME,
		POSTGRES_DB,
	)

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

func Init(test bool) {
	if test {
		POSTGRES_USER = os.Getenv("POSTGRES_TEST_USER")
		POSTGRES_PASSWORD = os.Getenv("POSTGRES_TEST_PASSWORD")
		POSTGRES_NAME = os.Getenv("POSTGRES_TEST_NAME")
		POSTGRES_DB = os.Getenv("POSTGRES_TEST_DB")
		REPORTS_DIRNAME = os.Getenv("REPORTS_TEST_DIRNAME")
	}

	connectToDB()
	ensureDir(filepath.Join(GENERATED_DIRNAME, REPORTS_DIRNAME))
}
