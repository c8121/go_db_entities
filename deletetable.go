package go_db_entities

import (
	"log"
	"reflect"
)

// DeleteQueryBuilder generates SQL-Statements.
type DeleteQueryBuilder interface {
	GetDeleteQuery(e Entity) string
}

// Delete takes an Entity and deletes it from the connected database.
func Delete(e Entity) error {

	if e == nil {
		return ErrNil
	}

	if e.GetId() == 0 {
		return ErrNotFound
	}

	sql := deleteBuilder.GetDeleteQuery(e)
	log.Println(sql)

	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Printf("Error preparing statement: %s", err)
		return err
	}
	defer CloseOrLog(stmt)

	v := reflect.ValueOf(GetAsStruct(e))

	values := make([]any, 1)
	values[0] = v.FieldByName("Id").Interface()
	log.Println(values)

	_, err = stmt.Exec(values...)
	if err != nil {
		log.Printf("Error executing query: %s", err)
		return err
	}

	e.SetId(0)

	return nil
}
