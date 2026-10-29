package sqlite

import (
	"fmt"
	"log/slog"

	"github.com/c8121/go_db_entities"
)

type SelectQueryBuilder struct {
}

func (b SelectQueryBuilder) GetSelectQuery(c *go_db_entities.Context, e go_db_entities.Entity, filter go_db_entities.Filter) string {

	sql := "SELECT "
	if filter.Distinct {
		sql += "DISTINCT "
	}

	sql += b.GetFields(c, e)

	sql += " FROM " + e.GetEntityName()

	where := c.QueryBuilder.FilterQueryBuilder.GetFilterStatement(c, filter)
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

func (b SelectQueryBuilder) GetFields(c *go_db_entities.Context, e go_db_entities.Entity) string {

	fields := ""
	properties := go_db_entities.GetProperties(c, e)
	entityFields := go_db_entities.GetVisibleFields(e)
	for _, field := range entityFields {
		if len(fields) > 0 {
			fields += ", "
		}

		fieldProps, hasFieldProps := properties[field.Name]

		typeName := field.Type.Name()
		if hasFieldProps && len(fieldProps.Type) > 0 {
			typeName = fieldProps.Type
		}

		typeMapping := c.GetTypeMapping(typeName)
		fields += typeMapping.GetFieldName(field.Name)
	}

	return fields
}
