package go_db_entities

import (
	"log"
	"reflect"
)

// UpdateQueryBuilder generates SQL-Statements.
type UpdateQueryBuilder interface {
	GetUpdateQuery(c *Context, e Entity) string
}

// Update takes an Entity and updates it into the connected database.
func Update(c *Context, e Entity) error {

	if e == nil {
		return ErrNil
	}

	if e.GetId() == 0 {
		return ErrNotFound
	}

	sql := c.QueryBuilder.UpdateQueryBuilder.GetUpdateQuery(c, e)
	log.Println(sql)

	stmt, err := c.DB.Prepare(sql)
	if err != nil {
		log.Printf("Error preparing statement: %s", err)
		return err
	}
	defer CloseOrLog(stmt)

	properties := GetProperties(c, e)
	entityFields := GetVisibleFields(e)
	v := reflect.ValueOf(GetAsStruct(e))

	values := make([]any, len(entityFields))
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
	values[p] = e.GetId()
	log.Println(values)

	_, err = stmt.Exec(values...)
	if err != nil {
		log.Printf("Error executing query: %s", err)
		return err
	}

	return nil
}
