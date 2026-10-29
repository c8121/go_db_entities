package sqlite

import (
	"github.com/c8121/go_db_entities"
)

type InsertQueryBuilder struct {
}

func (b InsertQueryBuilder) GetInsertQuery(c *go_db_entities.Context, e go_db_entities.Entity) string {

	sql := "INSERT INTO " + e.GetEntityName()

	fields := ""
	params := ""
	properties := go_db_entities.GetProperties(c, e)
	entityFields := go_db_entities.GetVisibleFields(e)
	for _, field := range entityFields {
		if field.Name != "Id" {
			if len(fields) > 0 {
				fields += ", "
				params += ", "
			}
			fields += field.Name

			fieldProps, hasFieldProps := properties[field.Name]

			typeName := field.Type.Name()
			if hasFieldProps && len(fieldProps.Type) > 0 {
				typeName = fieldProps.Type
			}

			typeMapping := c.GetTypeMapping(typeName)
			params += typeMapping.GetFieldParameter()
		}
	}

	sql += "(" + fields + ") VALUES (" + params + ");"

	return sql
}
