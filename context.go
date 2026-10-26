package go_db_entities

import "database/sql"

type Context struct {
	DB           *sql.DB
	QueryBuilder *QueryBuilder
}

func (c *Context) Close() error {
	return c.DB.Close()
}

// NewContext takes database connection and DDL/SQL-Builders required to read/write database.
func NewContext(database *sql.DB, queryBuilder *QueryBuilder) *Context {
	return &Context{
		DB:           database,
		QueryBuilder: queryBuilder,
	}
}
