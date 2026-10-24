package go_db_entities

import "errors"

var (
	ErrNotFound              = errors.New("Entity not found")
	ErrNotInsertable         = errors.New("not a Insertable")
	ErrNil                   = errors.New("nil object")
	ErrMustBePointerToStruct = errors.New("must be pointer to struct")
)
