package go_db_entities

import (
	"log/slog"
)

// DeleteQueryBuilder generates SQL-Statements.
type DeleteQueryBuilder interface {
	GetDeleteQuery(c *Storage, e Entity) string
	GetDeleteBatchQuery(c *Storage, e Entity, f Filter) string
}

// Delete takes an Entity and deletes it from the connected database.
func Delete(storage *Storage, e Entity) error {

	if e == nil {
		return ErrNil
	}

	if e.GetId() == 0 {
		return ErrNotFound
	}

	sql := storage.QueryBuilder.DeleteQueryBuilder.GetDeleteQuery(storage, e)
	slog.Debug(sql)

	stmt, err := storage.DB.Prepare(sql)
	if err != nil {
		slog.Error("Preparing statement failed", "error", err)
		return err
	}
	defer CloseOrLog(stmt)

	values := make([]any, 1)
	values[0] = e.GetId()
	slog.Debug("Using values", "v", values)

	_, err = stmt.Exec(values...)
	if err != nil {
		slog.Error("Executing query failed", "error", err)
		return err
	}

	e.SetId(0)

	return nil
}

// DeleteBatch takes an Entity and a Filter to delete multiple items from the connected database.
func DeleteBatch(storage *Storage, e Entity, f Filter) error {

	if e == nil {
		return ErrNil
	}

	sql := storage.QueryBuilder.DeleteQueryBuilder.GetDeleteBatchQuery(storage, e, f)
	slog.Debug(sql)

	stmt, err := storage.DB.Prepare(sql)
	if err != nil {
		slog.Error("Preparing statement failed", "error", err)
		return err
	}
	defer CloseOrLog(stmt)

	args := f.GetFilterValues()
	slog.Debug("Using args", "v", args)

	_, err = stmt.Exec(args...)
	if err != nil {
		slog.Error("Executing query failed", "error", err)
		return err
	}

	return nil
}
