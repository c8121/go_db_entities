package sqlite

import (
	"github.com/c8121/go_db_entities"
)

type DeleteQueryBuilder struct {
}

func (b DeleteQueryBuilder) GetDeleteQuery(e go_db_entities.Entity) string {

	sql := "DELETE FROM " + e.GetEntityName() + " WHERE Id = ?;"
	return sql
}
