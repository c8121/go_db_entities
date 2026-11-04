package go_db_entities

import (
	"database/sql"
)

type Context struct {
	DB            *sql.DB
	QueryBuilder  *QueryBuilder
	FieldMappings map[string]FieldMapping //Field-name -> TypeMapping
}

func (c *Context) GetTypeMapping(typeName string) FieldMapping {
	mapping, ok := c.FieldMappings[typeName]
	if !ok {
		return &DefaultFieldMapping{}
	}
	return mapping
}

// NewContext takes database connection and DDL/SQL-Builders required to read/write database.
func NewContext(database *sql.DB, queryBuilder *QueryBuilder, typeMapping map[string]FieldMapping) *Context {

	if database == nil {
		panic("database cannot be nil")
	}

	if queryBuilder == nil {
		panic("queryBuilder cannot be nil")
	}

	return &Context{
		DB:            database,
		QueryBuilder:  queryBuilder,
		FieldMappings: typeMapping,
	}
}
