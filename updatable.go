package go_db_entities

import (
	"log"
	"reflect"
)

// UpdateQueryBuilder generates SQL-Statements.
type UpdateQueryBuilder interface {
	GetUpdateQuery(e Entity) string
}

// Update takes an Entity and updates it into the connected database.
func Update(e Entity) error {

	if e == nil {
		return ErrNil
	}

	if e.GetId() == 0 {
		return ErrNotFound
	}

	sql := updateBuilder.GetUpdateQuery(e)
	log.Println(sql)

	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Printf("Error preparing statement: %s", err)
		return err
	}
	defer CloseOrLog(stmt)

	entityFields := GetVisibleFields(e)
	v := reflect.ValueOf(GetAsStruct(e))

	values := make([]any, len(entityFields))
	var p = 0
	for _, field := range entityFields {
		if field.Name != "Id" {
			values[p] = v.FieldByName(field.Name).Interface()
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
