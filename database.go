package go_db_entities

import "database/sql"

var (
	db            *sql.DB
	createBuilder CreateQueryBuilder
	insertBuilder InsertQueryBuilder
	updateBuilder UpdateQueryBuilder
	deleteBuilder DeleteQueryBuilder
	selectBuilder SelectQueryBuilder
)

func InitEntityDB(database *sql.DB,
	createQueryBuilder CreateQueryBuilder,
	insertQueryBuilder InsertQueryBuilder,
	updateQueryBuilder UpdateQueryBuilder,
	deleteQueryBuilder DeleteQueryBuilder,
	selectQueryBuilder SelectQueryBuilder,
) {
	db = database
	createBuilder = createQueryBuilder
	insertBuilder = insertQueryBuilder
	updateBuilder = updateQueryBuilder
	deleteBuilder = deleteQueryBuilder
	selectBuilder = selectQueryBuilder
}
