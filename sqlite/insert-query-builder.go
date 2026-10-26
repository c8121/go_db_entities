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
	entityFields := go_db_entities.GetVisibleFields(e)
	for _, field := range entityFields {
		if field.Name != "Id" {
			if len(fields) > 0 {
				fields += ", "
				params += ", "
			}
			fields += field.Name
			params += "?"
		}
	}

	sql += "(" + fields + ") VALUES (" + params + ");"

	return sql
}
