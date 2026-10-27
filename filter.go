package go_db_entities

import (
	"log"
	"reflect"
)

const (
	CompareEqual       = 1
	CompareGreater     = 2
	CompareGraterEqual = 3
	CompareLess        = 4
	CompareLessEqual   = 5
	CompareContains    = 6
	CompareIn          = 7

	ConjunctAnd = 1
	ConjunctOr  = 2
)

// FilterQueryBuilder generates SQL-WHERE-Statements.
type FilterQueryBuilder interface {
	GetFilterStatement(c *Context, f Filter) string
}

type Filter struct {
	FieldName  string
	Comparison int
	Value      any

	Next            *Filter
	NextConjunction int

	Offset int
	Count  int

	OrderBy []string
}

func NewFilter(fieldName string, comparison int, value any) *Filter {
	return &Filter{
		FieldName:  fieldName,
		Comparison: comparison,
		Value:      value,
	}
}

func NewIdFilter(id int64) *Filter {
	return NewFilter("Id", CompareEqual, id)
}

func (f *Filter) GetFilterValues() []any {

	values := make([]any, 0)

	if IsSlice(f.Value) {

		if f.Comparison != 0 && f.Comparison != CompareIn && f.Comparison != CompareEqual {
			log.Printf("WARN: Ignore invalid comparison for slice-field %s", f.FieldName)
		}

		v := reflect.ValueOf(f.Value)
		for i := 0; i < v.Len(); i++ {
			values = append(values, v.Index(i).Interface())
		}

	} else {

		values = append(values, f.Value)
		if f.Next != nil {
			values = append(values, f.Next.GetFilterValues()...)
		}

	}

	return values

}
