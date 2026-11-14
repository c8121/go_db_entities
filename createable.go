package go_db_entities

import (
	"log/slog"
)

// CreateQueryBuilder generates DDL-Statements to create or update tables.
type CreateQueryBuilder interface {
	// GetCreateQueries generates all statements to create a new table incl. indexes
	GetCreateQueries(session *Session, e Entity) []string
	// GetAddFieldsQueries generated statements to create fields in an existing table.
	GetAddFieldsQueries(session *Session, e Entity, fieldNames []string) []string
}

// Create takes an Entity and creates the table in the connected database.
// Uses CreateQueryBuilder to generate the DDL-Statement.
func Create(session *Session, e Entity) error {

	if e == nil {
		return ErrNil
	}

	ddl := session.Storage.QueryBuilder.CreateQueryBuilder.GetCreateQueries(session, e)
	for _, query := range ddl {
		slog.Debug(query)

		stmt, err := session.StatementProvider.Prepare(query)
		if err != nil {
			return err
		}

		_, err = stmt.Exec()
		if err != nil {
			CloseOrLog(stmt)
			slog.Error("DB.Exec failed", "error", err)
			return err
		}
		CloseOrLog(stmt)
	}

	return nil
}

// Validate takes an Entity and checks if all fields exist in the table of the connected database.
// Creates missing fields.
func Validate(session *Session, e Entity) error {

	if e == nil {
		return ErrNil
	}

	var missingFields []string
	checkedFields := CheckFieldsExist(session, e)
	for _, checkedField := range checkedFields {
		if checkedField.Error != nil {
			missingFields = append(missingFields, checkedField.FieldName)
		}
	}

	ddl := session.Storage.QueryBuilder.CreateQueryBuilder.GetAddFieldsQueries(session, e, missingFields)
	for _, query := range ddl {
		slog.Debug(query)

		stmt, err := session.StatementProvider.Prepare(query)
		if err != nil {
			return err
		}

		_, err = stmt.Exec()
		if err != nil {
			CloseOrLog(stmt)
			slog.Error("DB.Exec failed", "error", err)
			return err
		}
		CloseOrLog(stmt)
	}

	return nil
}
