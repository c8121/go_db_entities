package sqlite

import (
	"github.com/c8121/go_db_entities"
)

type CreateQueryBuilder struct {
}

func (b CreateQueryBuilder) GetCreateQuery(c *go_db_entities.Context, e go_db_entities.Entity) string {

	sql := "CREATE TABLE " + e.GetEntityName()

	fields := ""
	entityFields := go_db_entities.GetVisibleFields(e)
	for _, field := range entityFields {
		if len(fields) > 0 {
			fields += ", "
		}
		fields += field.Name

		if field.Name == "Id" {
			fields += " integer PRIMARY KEY"
		} else {
			switch field.Type.Name() {
			case "int":
				fields += " Integer"
				break
			case "string":
				fields += " Text"
				break
			case "Time":
				fields += " DateTime"
				break
			default:
				panic("Unknown Type: " + field.Type.Name())
			}
		}
	}

	return sql + "(" + fields + ");"
}
