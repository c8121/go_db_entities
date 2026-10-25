package go_db_entities

import (
	"log"
	"reflect"
)

type InsertQueryBuilder interface {
	GetInsertQuery(e Entity) string
}


func Insert(e Entity) error {

	if e == nil {
		return ErrNil
	}

	if e.GetId() != 0 {
		return ErrExists
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

	e.SetId(id)

	return nil
}
