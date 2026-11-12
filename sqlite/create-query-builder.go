package sqlite

import (
	"fmt"
	"reflect"

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

func (b CreateQueryBuilder) GetCreateMissingFieldsQueries(c *go_db_entities.Context, e go_db_entities.Entity) []string {

	var ddl []string

	properties := go_db_entities.GetProperties(c, e)
	entityFields := go_db_entities.GetVisibleFields(e)
	entityFieldsMap := make(map[string]reflect.StructField, len(entityFields))
	for _, field := range entityFields {
		entityFieldsMap[field.Name] = field
	}

	checkedFields := go_db_entities.CheckFieldsExist(c, e)
	for _, checkedField := range checkedFields {
		if checkedField.Error != nil {

			field, ok := entityFieldsMap[checkedField.FielName]
			if !ok {
				//panic is intentionally here because I assume it makes no sense to use this lib for data it cannot handle
				panic(fmt.Sprintf("CheckFieldsExist failed for field %s", checkedField.FielName))
			}

			sql := "ALTER TABLE " + e.GetEntityName() +
				" ADD COLUMN " + b.CreateFieldQuery(field, properties) + ";"

			ddl = append(ddl, sql)
		}
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
		fields += b.CreateFieldQuery(field, properties)
	}

	return sql + "(" + fields + ");"
}

func (b CreateQueryBuilder) CreateFieldQuery(field reflect.StructField, properties *go_db_entities.EntityProperties) string {

	sql := properties.GetFieldName(field.Name)

	typeName := properties.GetFieldType(field.Name, field.Type.Name())
	size := properties.GetFieldSize(field.Name)

	if field.Name == "Id" {
		sql += " integer PRIMARY KEY"
	} else {
		switch typeName {
		case "int", "int8", "int16", "int32", "int64":
			sql += " Integer"
			if size > 0 {
				sql += fmt.Sprintf("(%d)", size)
			}

		case "string":
			sql += " Text"
			if size > 0 {
				sql += fmt.Sprintf("(%d)", size)
			}

		case "Time":
			sql += " DateTime"

		case "Embedding":
			sql += fmt.Sprintf(" float[%d]", size)

		default:
			//panic is intentionally here because I assume it makes no sense to use this lib for data it cannot handle
			panic("Unknown Type: \"" + typeName + "\", Field: " + field.Name + "(" + field.Type.Name() + ")")
		}
	}
	return sql
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
