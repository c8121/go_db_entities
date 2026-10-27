package go_db_entities

import (
	"io"
	"log"
	"reflect"
)

// CloseOrLog closes the given Closer and logs error, if one occurs.
func CloseOrLog(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Printf("Close error: %v\n", err)
	}
}

func IsSlice(v interface{}) bool {
	if v == nil {
		return false
	}
	return reflect.TypeOf(v).Kind() == reflect.Slice
}
