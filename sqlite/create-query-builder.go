package sqlite

import (
	"fmt"

	"github.com/c8121/go_db_entities"
)

type CreateQueryBuilder struct {
}

func (b CreateQueryBuilder) GetCreateQueries(c *go_db_entities.Context, e go_db_entities.Entity) []string {

	ddl := []string{
		b.GetCreateTableQuery(c, e),
	}

	ddl = append(ddl, b.GetCreateIndexQueries(c, e)...)

	return ddl
}

func (b CreateQueryBuilder) GetCreateTableQuery(c *go_db_entities.Context, e go_db_entities.Entity) string {

	sql := "CREATE TABLE IF NOT EXISTS " + e.GetEntityName()

	fields := ""
	properties := go_db_entities.GetProperties(c, e)
	entityFields := go_db_entities.GetVisibleFields(e)
	for _, field := range entityFields {

		if len(fields) > 0 {
			fields += ", "
		}

		fields += properties.GetFieldName(field.Name)

		typeName := properties.GetFieldType(field.Name, field.Type.Name())
		size := properties.GetFieldSize(field.Name)

		if field.Name == "Id" {
			fields += " integer PRIMARY KEY"
		} else {
			switch typeName {
			case "int", "int8", "int16", "int32", "int64":
				fields += " Integer"
				if size > 0 {
					fields += fmt.Sprintf("(%d)", size)
				}

			case "string":
				fields += " Text"
				if size > 0 {
					fields += fmt.Sprintf("(%d)", size)
				}

			case "Time":
				fields += " DateTime"

			case "Embedding":
				fields += fmt.Sprintf(" float[%d]", size)

			default:
				panic("Unknown Type: \"" + typeName + "\", Field: " + field.Name + "(" + field.Type.Name() + ")")
			}
		}
	}

	return sql + "(" + fields + ");"
}

func (b CreateQueryBuilder) GetCreateIndexQueries(c *go_db_entities.Context, e go_db_entities.Entity) []string {

	queries := make([]string, 0)

	properties := go_db_entities.GetProperties(c, e)
	entityFields := go_db_entities.GetVisibleFields(e)
	for _, field := range entityFields {
		if properties.GetFieldIndex(field.Name) {

			fieldName := properties.GetFieldName(field.Name)

			queries = append(queries, "CREATE INDEX IF NOT EXISTS "+
				"IDX_"+e.GetEntityName()+"_"+fieldName+" "+
				"ON "+e.GetEntityName()+"("+fieldName+")")
		}
	}

	return queries
}
