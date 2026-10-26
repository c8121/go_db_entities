package go_db_entities

import "errors"

var (
	ErrNotFound              = errors.New("entity not found")
	ErrExists                = errors.New("entity already exists")
	ErrNil                   = errors.New("nil object")
	ErrMustBePointerToStruct = errors.New("must be pointer to struct")
)
