package sqlite

import (
	"github.com/c8121/go_db_entities"
)

type InsertQueryBuilder struct {
}

func (b InsertQueryBuilder) GetInsertQuery(session *go_db_entities.Session, e go_db_entities.Entity) string {

	sql := "INSERT INTO " + e.GetEntityName()

	fields := ""
	params := ""
	properties := go_db_entities.GetProperties(session, e)
	entityFields := go_db_entities.GetVisibleFields(e)
	for _, field := range entityFields {
		if field.Name != "Id" {
			if len(fields) > 0 {
				fields += ", "
				params += ", "
			}
			fields += properties.GetFieldName(field.Name)

			typeName := properties.GetFieldType(field.Name, field.Type.Name())

			typeMapping := session.Storage.GetTypeMapping(typeName)
			params += typeMapping.GetFieldParameter()
		}
	}

	sql += "(" + fields + ") VALUES (" + params + ");"

	return sql
}
