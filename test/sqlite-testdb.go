package test

import (
	"database/sql"
	"log"

	"github.com/c8121/go_db_entities"
	"github.com/c8121/go_db_entities/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

func OpenTestDB() *sql.DB {

	url := "file:TestDb?mode=memory&cache=shared"
	log.Printf("OpenTestDB: %s", url)

	db, err := sql.Open("sqlite3", url)
	if err != nil {
		panic(err)
	}

	go_db_entities.InitEntityDB(
		db,
		sqlite.CreateQueryBuilder{},
		sqlite.InsertQueryBuilder{},
		sqlite.UpdateQueryBuilder{},
		sqlite.DeleteQueryBuilder{},
		sqlite.SelectQueryBuilder{},
	)

	return db
}
