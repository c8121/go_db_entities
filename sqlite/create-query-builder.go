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

	sql := "CREATE TABLE " + e.GetEntityName()

	fields := ""
	properties := go_db_entities.GetProperties(c, e)
	entityFields := go_db_entities.GetVisibleFields(e)
	for _, field := range entityFields {

		fieldProps, hasFieldProps := properties[field.Name]

		if len(fields) > 0 {
			fields += ", "
		}

		if hasFieldProps && len(fieldProps.Name) > 0 {
			fields += fieldProps.Name
		} else {
			fields += field.Name
		}

		typeName := field.Type.Name()
		if hasFieldProps && len(fieldProps.Type) > 0 {
			typeName = fieldProps.Type
		}

		if field.Name == "Id" {
			fields += " integer PRIMARY KEY"
		} else {
			switch typeName {
			case "int", "int8", "int16", "int32", "int64":
				fields += " Integer"
				if hasFieldProps && fieldProps.Size > 0 {
					fields += fmt.Sprintf("(%d)", fieldProps.Size)
				}

			case "string":
				fields += " Text"
				if hasFieldProps && fieldProps.Size > 0 {
					fields += fmt.Sprintf("(%d)", fieldProps.Size)
				}

			case "Time":
				fields += " DateTime"

			case "Embedding":
				fields += fmt.Sprintf(" float[%d]", fieldProps.Size)

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
		fieldProps, hasFieldProps := properties[field.Name]
		if hasFieldProps {
			if fieldProps.Index {
				queries = append(queries, "CREATE INDEX "+
					"IDX_"+e.GetEntityName()+"_"+field.Name+" "+
					"ON "+e.GetEntityName()+"("+field.Name+")")
			}
		}
	}

	return queries
}
