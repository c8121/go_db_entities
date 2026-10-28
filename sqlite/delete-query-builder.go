package sqlite

import (
	"github.com/c8121/go_db_entities"
)

type DeleteQueryBuilder struct {
}

func (b DeleteQueryBuilder) GetDeleteQuery(c *go_db_entities.Context, e go_db_entities.Entity) string {

	sql := "DELETE FROM " + e.GetEntityName() + " WHERE Id=?;"
	return sql
}

func (b DeleteQueryBuilder) GetDeleteBatchQuery(c *go_db_entities.Context, e go_db_entities.Entity, f go_db_entities.Filter) string {
	sql := "DELETE FROM " + e.GetEntityName() + " WHERE " + c.QueryBuilder.FilterQueryBuilder.GetFilterStatement(c, f) + ";"
	return sql
}
