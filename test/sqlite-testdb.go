package test

import (
	"database/sql"
	"log/slog"

	"github.com/c8121/go_db_entities"
	"github.com/c8121/go_db_entities/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

func OpenTestDB() *go_db_entities.Context {

	url := "file:TestDb?mode=memory&cache=shared"
	slog.Info("OpenTestDB", "url", url)

	db, err := sql.Open("sqlite3", url)
	if err != nil {
		panic(err)
	}

	return go_db_entities.NewContext(
		db,
		sqlite.NewQueryBuilder(),
		sqlite.NewFieldMapping(),
	)
}
