package sqlite

import (
	"fmt"

	"github.com/c8121/go_db_entities"
)

type TableQueryBuilder struct {
}

func (b TableQueryBuilder) GetTableQueryStatement(c *go_db_entities.Context, q go_db_entities.TableQuery) string {

	sql := "SELECT " + b.GetFields(c, q) + " " +
		"FROM " + b.GetFrom(c, q)

	where := b.GetWhere(c, q)
	if len(where) > 0 {
		sql += " WHERE " + where
	}

	limit := b.GetLimit(c, q)
	if len(limit) > 0 {
		sql += " LIMIT " + limit
	}

	sql += ";"

	return sql
}

func (b TableQueryBuilder) GetLimit(c *go_db_entities.Context, q go_db_entities.TableQuery) string {

	sql := ""

	if q.Filter != nil && (q.Filter.Offset > 0 || q.Filter.Count > 0) {
		sql = fmt.Sprintf("%d,%d", q.Filter.Offset, q.Filter.Count)
	}

	return sql
}

func (b TableQueryBuilder) GetWhere(c *go_db_entities.Context, q go_db_entities.TableQuery) string {

	sql := ""

	if q.Filter != nil {
		sql = c.QueryBuilder.FilterQueryBuilder.GetFilterStatement(c, *q.Filter)
	}

	if q.Join != nil {
		for _, join := range q.Join {
			if join.TableQuery.Filter != nil {
				if len(sql) > 0 {
					sql += " AND "
				}
				sql += c.QueryBuilder.FilterQueryBuilder.GetFilterStatement(c, *join.TableQuery.Filter)
			}
		}
	}

	return sql
}

func (b TableQueryBuilder) GetFrom(c *go_db_entities.Context, q go_db_entities.TableQuery) string {

	sql := q.Entity.GetEntityName()

	if q.Join != nil {
		for _, join := range q.Join {
			switch join.Type {
			case go_db_entities.InnerJoin:
				sql += " INNER JOIN "
			case go_db_entities.LeftJoin:
				sql += " LEFT JOIN "
			default:
				sql += " INNER JOIN "
			}
			sql += b.GetFrom(c, *join.TableQuery)
			sql += " ON " + q.Entity.GetEntityName() + "." + join.OnFrom + "=" + join.TableQuery.Entity.GetEntityName() + "." + join.OnTo
			sql = "(" + sql + ")"
		}
	}

	return sql
}

func (b TableQueryBuilder) GetFields(c *go_db_entities.Context, q go_db_entities.TableQuery) string {

	sql := ""
	for _, field := range q.FieldNames {
		if len(sql) > 0 {
			sql += ","
		}
		sql += q.Entity.GetEntityName() + "." + field
	}

	if q.Join != nil {
		for _, join := range q.Join {
			nextFields := b.GetFields(c, *join.TableQuery)
			if len(nextFields) > 0 {
				if len(sql) > 0 {
					sql += ","
				}
				sql += nextFields
			}
		}
	}

	return sql
}
