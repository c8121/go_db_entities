package go_db_entities

import (
	"database/sql"
)

type Context struct {
	DB            *sql.DB
	QueryBuilder  *QueryBuilder
	FieldMappings map[string]FieldMapping //Fieldname -> TypeMapping
}

func (c *Context) GetTypeMapping(typeName string) FieldMapping {
	mapping, ok := c.FieldMappings[typeName]
	if !ok {
		return &DefaultFieldMapping{}
	}
	return mapping
}

func (c *Context) Close() error {
	return c.DB.Close()
}

// NewContext takes database connection and DDL/SQL-Builders required to read/write database.
func NewContext(database *sql.DB, queryBuilder *QueryBuilder, typeMapping map[string]FieldMapping) *Context {
	return &Context{
		DB:            database,
		QueryBuilder:  queryBuilder,
		FieldMappings: typeMapping,
	}
}
