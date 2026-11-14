package go_db_entities

import (
	"database/sql"
)

type Storage struct {
	DB            *sql.DB
	QueryBuilder  *QueryBuilder
	FieldMappings map[string]FieldMapping //Field-name -> TypeMapping
}

// NewStorage takes database connection and DDL/SQL-Builders required to read/write database.
func NewStorage(database *sql.DB, queryBuilder *QueryBuilder, typeMapping map[string]FieldMapping) *Storage {

	if database == nil {
		panic("database cannot be nil")
	}

	if queryBuilder == nil {
		panic("queryBuilder cannot be nil")
	}

	return &Storage{
		DB:            database,
		QueryBuilder:  queryBuilder,
		FieldMappings: typeMapping,
	}
}
