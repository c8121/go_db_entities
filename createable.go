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
	// GetAddFieldsQueries generated statements to create fields in an existing table.
	GetAddFieldsQueries(c *Storage, e Entity, fieldNames []string) []string
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

	var missingFields []string
	checkedFields := CheckFieldsExist(c, e)
	for _, checkedField := range checkedFields {
		if checkedField.Error != nil {
			missingFields = append(missingFields, checkedField.FielName)
		}
	}

	ddl := c.QueryBuilder.CreateQueryBuilder.GetAddFieldsQueries(c, e, missingFields)
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
