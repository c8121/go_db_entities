package go_db_entities

const (
	CompareEqual       = 1
	CompareGreater     = 2
	CompareGraterEqual = 3
	CompareLess        = 4
	CompareLessEqual   = 5
	CompareContains    = 6

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

	values = append(values, f.Value)
	if f.Next != nil {
		values = append(values, f.Next.GetFilterValues()...)
	}

	return values

}
