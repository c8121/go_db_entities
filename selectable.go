package go_db_entities

import (
	"log"
	"reflect"
)

// SelectQueryBuilder generates SQL-Statements.
type SelectQueryBuilder interface {
	GetSelectQuery(e Entity, filter map[string]any) string
}

// Select takes a filter-map (column-name -> value), selects from the connected database and assigns to given Entity.
func Select(e Entity, filter map[string]any) error {

	if e == nil {
		return ErrNil
	}

	destValue := reflect.ValueOf(e)
	if destValue.Kind() != reflect.Ptr || destValue.Elem().Kind() != reflect.Struct {
		return ErrMustBePointerToStruct
	}

	sql := selectBuilder.GetSelectQuery(e, filter)
	log.Println(sql)

	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Printf("Error preparing statement: %s", err)
		return err
	}
	defer CloseOrLog(stmt)

	args := make([]any, 0)
	for _, arg := range filter {
		args = append(args, arg)
	}
	log.Println(args)

	rows, err := stmt.Query(args...)
	if err != nil {
		log.Printf("Error executing query: %s", err)
		return err
	}
	defer CloseOrLog(rows)

	if rows.Next() {

		structValue := destValue.Elem()
		entityFields := GetVisibleFields(e)

		values := make([]any, len(entityFields))
		scanArgs := make([]interface{}, len(values))

		for i, field := range entityFields {
			destField := structValue.FieldByName(field.Name)
			scanArgs[i] = destField.Addr().Interface()
		}

		err = rows.Scan(scanArgs...)
		if err != nil {
			return err
		}
	} else {
		e.SetId(0)
		return ErrNotFound
	}

	return nil
}
