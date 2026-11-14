package sqlite

import (
	"github.com/c8121/go_db_entities"
)

type DeleteQueryBuilder struct {
}

func (b DeleteQueryBuilder) GetDeleteQuery(storage *go_db_entities.Storage, e go_db_entities.Entity) string {

	sql := "DELETE FROM " + e.GetEntityName() + " WHERE Id=?;"
	return sql
}

func (b DeleteQueryBuilder) GetDeleteBatchQuery(storage *go_db_entities.Storage, e go_db_entities.Entity, f go_db_entities.Filter) string {
	sql := "DELETE FROM " + e.GetEntityName() + " WHERE " + storage.QueryBuilder.FilterQueryBuilder.GetFilterStatement(storage, f) + ";"
	return sql
}
