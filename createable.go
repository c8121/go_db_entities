package go_db_entities

import (
	"log/slog"
)

const (
	StructTagName = "dbe"
)

// CreateQueryBuilder generates DDL-Statements to create or update tables.
type CreateQueryBuilder interface {
	// GetCreateQueries generates all statements to create a new table incl. indexes
	GetCreateQueries(c *Storage, e Entity) []string
	// GetCreateMissingFieldsQueries generated only those statements that are required to create fields which are missing in an existing table.
	GetCreateMissingFieldsQueries(c *Storage, e Entity) []string
}

// Create takes an Entity and creates the table in the connected database.
// Uses CreateQueryBuilder to generate the DDL-Statement.
func Create(c *Storage, e Entity) error {

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

// Validate takes an Entity and checks if all fields exist in the table of the connected database.
// Creates missing fields.
func Validate(c *Storage, e Entity) error {

	if e == nil {
		return ErrNil
	}

	ddl := c.QueryBuilder.CreateQueryBuilder.GetCreateMissingFieldsQueries(c, e)
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
