package go_db_entities

import (
	"reflect"
)

// Entity provides function required to build SQL/DDL-Statements
// and to identify a database entity.
type Entity interface {
	GetEntityName() string

	GetId() int64
	SetId(id int64)
}

// GetAsStruct takes an interface (given as parameter to a function for example an unwraps the struct.
func GetAsStruct(e Entity) any {
	ptr := reflect.ValueOf(e)
	wrp := reflect.Indirect(ptr)
	return wrp.Interface()
}

// GetVisibleFields reads all visible fields from the given Entity.
func GetVisibleFields(e Entity) []reflect.StructField {

	value := GetAsStruct(e)
	return reflect.VisibleFields(reflect.TypeOf(value))
}
