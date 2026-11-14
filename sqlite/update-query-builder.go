package sqlite

import (
	"github.com/c8121/go_db_entities"
)

type UpdateQueryBuilder struct {
}

func (b UpdateQueryBuilder) GetUpdateQuery(session *go_db_entities.Session, e go_db_entities.Entity) string {

	sql := "UPDATE " + e.GetEntityName()

	fields := ""
	properties := go_db_entities.GetProperties(session, e)
	entityFields := go_db_entities.GetVisibleFields(e)
	for _, field := range entityFields {
		if field.Name != "Id" {
			if len(fields) > 0 {
				fields += ", "
			}

			fields += properties.GetFieldName(field.Name) + "="

			typeName := properties.GetFieldType(field.Name, field.Type.Name())

			typeMapping := session.Storage.GetTypeMapping(typeName)
			fields += typeMapping.GetFieldParameter()
		}
	}

	sql += " SET " + fields + " WHERE Id=?;"

	return sql
}
