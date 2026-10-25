package go_db_entities

import (
	"reflect"
)

type Entity interface {
	GetEntityName() string

	GetId() int64
	SetId(id int64)
}

func GetAsStruct(e Entity) any {
	ptr := reflect.ValueOf(e)
	wrp := reflect.Indirect(ptr)
	return wrp.Interface()
}

func GetVisibleFields(e Entity) []reflect.StructField {

	ptr := reflect.ValueOf(e)
	wrp := reflect.Indirect(ptr)
	value := wrp.Interface()

	return reflect.VisibleFields(reflect.TypeOf(value))
}
