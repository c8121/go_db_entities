package go_db_entities

import (
	"log/slog"
	"reflect"
)

// InsertQueryBuilder generates SQL-Statements.
type InsertQueryBuilder interface {
	GetInsertQuery(c *Context, e Entity) string
}

// Insert takes an Entity and inserts it into the connected database.
func Insert(c *Context, e Entity) error {

	if e == nil {
		return ErrNil
	}

	if e.GetId() != 0 {
		return ErrExists
	}

	sql := c.QueryBuilder.InsertQueryBuilder.GetInsertQuery(c, e)
	slog.Debug(sql)

	stmt, err := c.DB.Prepare(sql)
	if err != nil {
		slog.Error("Preparing statement failed", "error", err)
		return err
	}
	defer CloseOrLog(stmt)

	properties := GetProperties(c, e)
	entityFields := GetVisibleFields(e)
	v := reflect.ValueOf(GetAsStruct(e))

	values := make([]any, len(entityFields)-1)
	var p = 0
	for _, field := range entityFields {
		if field.Name != "Id" {

			fieldProps, hasFieldProps := properties[field.Name]

			typeName := field.Type.Name()
			if hasFieldProps && len(fieldProps.Type) > 0 {
				typeName = fieldProps.Type
			}

			values[p] = v.FieldByName(field.Name).Interface()

			typeMapping := c.GetTypeMapping(typeName)
			values[p], err = typeMapping.ConvertStructValue(values[p])
			if err != nil {
				return err
			}

			p++
		}
	}
	slog.Debug("Using values", "v", values)

	res, err := stmt.Exec(values...)
	if err != nil {
		slog.Error("Executing query failed", "error", err)
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		slog.Error("Retrieving id failed", "error", err)
		return err
	}

	e.SetId(id)

	return nil
}
