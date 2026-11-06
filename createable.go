package go_db_entities

import (
	"log/slog"
)

const (
	StructTagName = "dbe"
)

// CreateQueryBuilder generates DDL-Statements.
type CreateQueryBuilder interface {
	GetCreateQueries(c *Context, e Entity) []string
}

// Create takes an Entity and creates the table in the connected database.
// Uses CreateQueryBuilder to generate the DDL-Statement.
func Create(c *Context, e Entity) error {

	if e == nil {
		return ErrNil
	}

	ddl := c.QueryBuilder.CreateQueryBuilder.GetCreateQueries(c, e)
	for _, query := range ddl {
		slog.Debug(query)

		_, err := c.DB.Exec(query)
		if err != nil {
			slog.Error("DB.Exec failed", "error", err)
			return err
		}
	}

	return nil
}
