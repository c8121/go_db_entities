package go_db_entities

import "errors"

var (
	ErrNotFound              = errors.New("Entity not found")
	ErrExists                = errors.New("Entity already exists")
	ErrNil                   = errors.New("nil object")
	ErrMustBePointerToStruct = errors.New("must be pointer to struct")
)
