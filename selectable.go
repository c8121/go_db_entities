package go_db_entities

import (
	"log/slog"
	"reflect"
)

// SelectQueryBuilder generates SQL-Statements.
type SelectQueryBuilder interface {
	GetSelectQuery(c *Context, e Entity, filter Filter) string
}

// Select takes a filter, selects rows from the connected database
// and calls entitySupplier for each row to apply values to Entity.
func Select(c *Context, entitySupplier func() Entity, filter Filter) ([]Entity, error) {

	e := entitySupplier()

	sql := c.QueryBuilder.SelectQueryBuilder.GetSelectQuery(c, e, filter)
	slog.Debug(sql)

	stmt, err := c.DB.Prepare(sql)
	if err != nil {
		slog.Error("Preparing statement failed", "error", err)
		return nil, err
	}
	defer CloseOrLog(stmt)

	args := filter.GetFilterValues()
	slog.Debug("Using args", "v", args)

	rows, err := stmt.Query(args...)
	if err != nil {
		slog.Error("Executing query failed", "error", err)
		return nil, err
	}
	defer CloseOrLog(rows)

	var entities = make([]Entity, 0)
	var properties map[string]FieldProperties = nil
	var entityFields []reflect.StructField = nil

	for rows.Next() {

		if len(entities) > 0 {
			//First Entity was created at top
			e = entitySupplier()
		}

		destValue := reflect.ValueOf(e)
		if destValue.Kind() != reflect.Ptr || destValue.Elem().Kind() != reflect.Struct {
			return nil, ErrMustBePointerToStruct
		}

		structValue := destValue.Elem()
		if entityFields == nil {
			properties = GetProperties(c, e)
			entityFields = GetVisibleFields(e)
		}
		scanArgs := make([]any, len(entityFields))
		scanValues := make([]any, len(entityFields))
		for i := 0; i < len(entityFields); i++ {
			scanArgs[i] = &scanValues[i]
		}

		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(entityFields); i++ {
			field := entityFields[i]
			destField := structValue.FieldByName(field.Name)

			fieldProps, hasFieldProps := properties[field.Name]
			typeName := field.Type.Name()
			if hasFieldProps && len(fieldProps.Type) > 0 {
				typeName = fieldProps.Type
			}

			typeMapping := c.GetTypeMapping(typeName)
			value, err := typeMapping.ConvertDbValue(scanValues[i])
			if err != nil {
				return nil, err
			}
			dbVal := reflect.ValueOf(value)
			destField.Set(dbVal.Convert(destField.Type()))
		}

		entities = append(entities, e)
	}

	if len(entities) == 0 {
		return entities, ErrNotFound
	}

	return entities, nil
}
