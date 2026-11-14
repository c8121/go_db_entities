package sqlite

import (
	"github.com/c8121/go_db_entities"
)

type DeleteQueryBuilder struct {
}

func (b DeleteQueryBuilder) GetDeleteQuery(session *go_db_entities.Session, e go_db_entities.Entity) string {

	sql := "DELETE FROM " + e.GetEntityName() + " WHERE Id=?;"
	return sql
}

func (b DeleteQueryBuilder) GetDeleteBatchQuery(session *go_db_entities.Session, e go_db_entities.Entity, f go_db_entities.Filter) string {
	sql := "DELETE FROM " + e.GetEntityName() + " WHERE " + session.Storage.QueryBuilder.FilterQueryBuilder.GetFilterStatement(session, f) + ";"
	return sql
}
