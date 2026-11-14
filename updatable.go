package go_db_entities

import (
	"log/slog"
	"reflect"
)

// UpdateQueryBuilder generates SQL-Statements.
type UpdateQueryBuilder interface {
	GetUpdateQuery(storage *Storage, e Entity) string
}

// Update takes an Entity and updates it into the connected database.
func Update(storage *Storage, e Entity) error {

	if e == nil {
		return ErrNil
	}

	if e.GetId() == 0 {
		return ErrNotFound
	}

	sql := storage.QueryBuilder.UpdateQueryBuilder.GetUpdateQuery(storage, e)
	slog.Debug(sql)

	stmt, err := storage.DB.Prepare(sql)
	if err != nil {
		slog.Error("Preparing statement failed", "error", err)
		return err
	}
	defer CloseOrLog(stmt)

	properties := GetProperties(storage, e)
	entityFields := GetVisibleFields(e)
	v := reflect.ValueOf(GetAsStruct(e))

	values := make([]any, len(entityFields))
	var p = 0
	for _, field := range entityFields {
		if field.Name != "Id" {

			typeName := properties.GetFieldType(field.Name, field.Type.Name())

			values[p] = v.FieldByName(field.Name).Interface()

			typeMapping := storage.GetTypeMapping(typeName)
			values[p], err = typeMapping.ConvertStructValue(values[p])
			if err != nil {
				return err
			}

			p++
		}
	}
	values[p] = e.GetId()
	slog.Debug("Using values", "v", values)

	res, err := stmt.Exec(values...)
	if err != nil {
		slog.Error("Executing query failed", "error", err)
		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		slog.Error("Getting affected rows failed", "error", err)
		return err
	}
	if cnt == 0 {
		return ErrNotFound
	}

	return nil
}
