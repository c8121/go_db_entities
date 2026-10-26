package sqlite

import "github.com/c8121/go_db_entities"

type SelectQueryBuilder struct {
}

func (b SelectQueryBuilder) GetSelectQuery(c *go_db_entities.Context, e go_db_entities.Entity, filter go_db_entities.Filter) string {

	sql := "SELECT * FROM " + e.GetEntityName() + " WHERE "
	sql += c.QueryBuilder.FilterQueryBuilder.GetFilterStatement(c, filter) + ";"
	return sql
}
