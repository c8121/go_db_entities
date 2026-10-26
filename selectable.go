package go_db_entities

import (
	"log"
	"reflect"
)

// SelectQueryBuilder generates SQL-Statements.
type SelectQueryBuilder interface {
	GetSelectQuery(c *Context, e Entity, filter Filter) string
}

// Select takes a filter-map (column-name -> value), selects rows from the connected database
// and calls entitySupplier for each row to apply values to Entity.
func Select(c *Context, entitySupplier func() Entity, filter Filter) ([]Entity, error) {

	e := entitySupplier()

	sql := c.QueryBuilder.SelectQueryBuilder.GetSelectQuery(c, e, filter)
	log.Println(sql)

	stmt, err := c.DB.Prepare(sql)
	if err != nil {
		log.Printf("Error preparing statement: %s", err)
		return nil, err
	}
	defer CloseOrLog(stmt)

	args := filter.GetFilterValues()
	log.Println(args)

	rows, err := stmt.Query(args...)
	if err != nil {
		log.Printf("Error executing query: %s", err)
		return nil, err
	}
	defer CloseOrLog(rows)

	var entities = make([]Entity, 0)
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
			entityFields = GetVisibleFields(e)
		}

		values := make([]any, len(entityFields))
		scanArgs := make([]interface{}, len(values))

		for i, field := range entityFields {
			destField := structValue.FieldByName(field.Name)
			scanArgs[i] = destField.Addr().Interface()
		}

		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		entities = append(entities, e)
	}

	if len(entities) == 0 {
		return entities, ErrNotFound
	}

	return entities, nil
}
