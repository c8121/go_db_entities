package go_db_entities

import "errors"

type CheckFieldResult struct {
	FieldName string
	Error     error
}

func CheckFieldsExist(session *Session, e Entity) []CheckFieldResult {

	var list []CheckFieldResult

	entityFields := GetVisibleFields(e)
	for _, field := range entityFields {

		query := &TableQuery{
			FieldNames: []string{field.Name},
			Entity:     e,
			Filter:     NewIdFilter(-1),
		}

		_, err := Query(session, *query)
		if errors.Is(err, ErrNotFound) {
			err = nil
		}
		list = append(list, CheckFieldResult{
			FieldName: field.Name,
			Error:     err,
		})
	}

	return list
}
