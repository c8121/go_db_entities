package test

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func OpenTestDB() *sql.DB {

	url := "file:TestDb?mode=memory&cache=shared"
	log.Printf("OpenTestDB: %s", url)

	db, err := sql.Open("sqlite3", url)
	if err != nil {
		panic(err)
	}

	return db
}
