package sqlite

import (
	"fmt"
	"log/slog"

	"github.com/c8121/go_db_entities"
)

type SelectQueryBuilder struct {
}

func (b SelectQueryBuilder) GetSelectQuery(session *go_db_entities.Session, e go_db_entities.Entity, filter go_db_entities.Filter) string {

	sql := "SELECT "
	if filter.Distinct {
		sql += "DISTINCT "
	}

	sql += b.GetFields(session, e)

	sql += " FROM " + e.GetEntityName()

	where := session.Storage.QueryBuilder.FilterQueryBuilder.GetFilterStatement(session, filter)
	if len(where) > 0 {
		sql += " WHERE " + where
	}

	if filter.OrderBy != nil {
		order := ""
		for _, field := range filter.OrderBy {
			if len(order) > 0 {
				order += ","
			}
			order += field
		}
		if len(order) > 0 {
			sql += " ORDER BY " + order
		}
	}

	if filter.Offset > 0 || filter.Count > 0 {
		sql += fmt.Sprintf(" LIMIT %d,%d", filter.Offset, filter.Count)
	}

	if filter.Next != nil && filter.Next.Offset > 0 && filter.Next.Count > 0 {
		slog.Warn("Ignoring offest/count of nested filter")
	}

	sql += ";"
	return sql
}

func (b SelectQueryBuilder) GetFields(session *go_db_entities.Session, e go_db_entities.Entity) string {

	fields := ""
	properties := go_db_entities.GetProperties(session, e)
	entityFields := go_db_entities.GetVisibleFields(e)
	for _, field := range entityFields {
		if len(fields) > 0 {
			fields += ", "
		}

		fieldName := properties.GetFieldName(field.Name)
		typeName := properties.GetFieldType(field.Name, field.Type.Name())

		typeMapping := session.Storage.GetTypeMapping(typeName)
		fields += typeMapping.GetFieldName(fieldName)
	}

	return fields
}
