package sqlite

import (
	"github.com/c8121/go_db_entities"
)

type UpdateQueryBuilder struct {
}

func (b UpdateQueryBuilder) GetUpdateQuery(c *go_db_entities.Context, e go_db_entities.Entity) string {

	sql := "UPDATE " + e.GetEntityName()

	fields := ""
	entityFields := go_db_entities.GetVisibleFields(e)
	for _, field := range entityFields {
		if field.Name != "Id" {
			if len(fields) > 0 {
				fields += ", "
			}
			fields += field.Name + "=?"
		}
	}

	sql += " SET " + fields + " WHERE Id=?;"

	return sql
}
