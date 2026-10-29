package go_db_entities

type FieldMapping interface {

	// GetFieldName returns the field name to be used in select statement.
	GetFieldName(name string) string

	// GetFieldParameter returns the SQL-Parameter used for insert/update statements.
	// Usually '?'
	GetFieldParameter() string

	// ConvertStructValue takes the value of a struct field and converts it to be used for insert/update statements.
	ConvertStructValue(value any) (any, error)

	// ConvertDbValue takes the value out of the database and converts it to be assigned to struct-field.
	ConvertDbValue(value any) (any, error)
}

type DefaultFieldMapping struct {
}

func (m *DefaultFieldMapping) GetFieldName(name string) string {
	return name
}

func (m *DefaultFieldMapping) GetFieldParameter() string {
	return "?"
}

func (m *DefaultFieldMapping) ConvertStructValue(value any) (any, error) {
	return value, nil
}

func (m *DefaultFieldMapping) ConvertDbValue(value any) (any, error) {
	return value, nil
}
