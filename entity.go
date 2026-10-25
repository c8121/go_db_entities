package go_db_entities

import (
	"reflect"
)

// Entity provides functions required to build SQL/DDL-Statements
// and to identify a database entity.
type Entity interface {
	GetEntityName() string

	GetId() int64
	SetId(id int64)
}

// Load selects from connected database and applies to given Entity
func Load(e Entity, filter map[string]any) error {

	if e == nil {
		return ErrNil
	}

	err := Select(e, filter)
	if err != nil {
		return err
	}

	return nil
}

// Save checks if Entity has an Id and either inserts or updates it in connected database.
func Save(e Entity) error {

	if e == nil {
		return ErrNil
	}

	if e.GetId() == 0 {
		return Insert(e)
	}

	return Update(e)
}

// CreateIdFilter creates a map which contains Id = id
func CreateIdFilter(id int64) map[string]any {
	filter := make(map[string]any)
	filter["Id"] = id
	return filter
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
