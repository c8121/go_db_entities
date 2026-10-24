package go_db_entities

import (
	"log"
	"reflect"
)

type InsertQueryBuilder interface {
	GetInsertQuery(e Entity) string
}

type Insertable interface {
	SetId(id int64)
}

func Insert(e Entity) error {

	if e == nil {
		return ErrNil
	}

	insertable, ok := e.(Insertable)
	if !ok {
		return ErrNotInsertable
	}

	sql := insertBuilder.GetInsertQuery(e)
	log.Println(sql)

	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Printf("Error preparing statement: %s", err)
		return err
	}
	defer CloseOrLog(stmt)

	entityFields := GetVisibleFields(e)
	v := reflect.ValueOf(GetAsStruct(e))

	values := make([]any, len(entityFields)-1)
	var p = 0
	for _, field := range entityFields {
		if field.Name != "Id" {
			values[p] = v.FieldByName(field.Name).Interface()
			p++
		}
	}
	log.Println(values)

	res, err := stmt.Exec(values...)
	if err != nil {
		log.Printf("Error executing query: %s", err)
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Error retrieving id: %s", err)
		return err
	}

	insertable.SetId(id)

	return nil
}
