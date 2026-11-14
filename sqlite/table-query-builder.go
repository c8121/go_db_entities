package sqlite

import (
	"fmt"

	"github.com/c8121/go_db_entities"
)

type TableQueryBuilder struct {
}

func (b TableQueryBuilder) GetTableQueryStatement(session *go_db_entities.Session, q go_db_entities.TableQuery) string {

	sql := "SELECT "
	if q.Filter != nil && q.Filter.Distinct {
		sql += "DISTINCT "
	}

	sql += b.GetFields(session, q) + " " +
		"FROM " + b.GetFrom(session, q)

	where := b.GetWhere(session, q)
	if len(where) > 0 {
		sql += " WHERE " + where
	}

	order := b.GetOrderBy(session, q)
	if len(order) > 0 {
		sql += " ORDER BY " + order
	}

	limit := b.GetLimit(session, q)
	if len(limit) > 0 {
		sql += " LIMIT " + limit
	}

	sql += ";"

	return sql
}

func (b TableQueryBuilder) GetOrderBy(session *go_db_entities.Session, q go_db_entities.TableQuery) string {

	sql := ""

	if q.Filter != nil && q.Filter.OrderBy != nil {
		for _, field := range q.Filter.OrderBy {
			if len(sql) > 0 {
				sql += ","
			}
			sql += field
		}
	}

	return sql
}

func (b TableQueryBuilder) GetLimit(session *go_db_entities.Session, q go_db_entities.TableQuery) string {

	sql := ""

	if q.Filter != nil && (q.Filter.Offset > 0 || q.Filter.Count > 0) {
		sql = fmt.Sprintf("%d,%d", q.Filter.Offset, q.Filter.Count)
	}

	return sql
}

func (b TableQueryBuilder) GetWhere(session *go_db_entities.Session, q go_db_entities.TableQuery) string {

	sql := ""

	if q.Filter != nil {
		sql = session.Storage.QueryBuilder.FilterQueryBuilder.GetFilterStatement(session, *q.Filter)
	}

	if q.Join != nil {
		for _, join := range q.Join {
			if join.TableQuery.Filter != nil {
				if len(sql) > 0 {
					sql += " AND "
				}
				sql += session.Storage.QueryBuilder.FilterQueryBuilder.GetFilterStatement(session, *join.TableQuery.Filter)
			}
		}
	}

	return sql
}

func (b TableQueryBuilder) GetFrom(session *go_db_entities.Session, q go_db_entities.TableQuery) string {

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
			sql += b.GetFrom(session, *join.TableQuery)
			sql += " ON " + q.Entity.GetEntityName() + "." + join.OnFrom + "=" + join.TableQuery.Entity.GetEntityName() + "." + join.OnTo
			sql = "(" + sql + ")"
		}
	}

	return sql
}

func (b TableQueryBuilder) GetFields(session *go_db_entities.Session, q go_db_entities.TableQuery) string {

	sql := ""
	for _, field := range q.FieldNames {
		if len(sql) > 0 {
			sql += ","
		}
		sql += q.Entity.GetEntityName() + "." + field
	}

	if q.Join != nil {
		for _, join := range q.Join {
			nextFields := b.GetFields(session, *join.TableQuery)
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
