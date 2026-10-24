package sqlite

import "github.com/c8121/go_db_entities"

type SelectQueryBuilder struct {
}

func (b SelectQueryBuilder) GetSelectQuery(e go_db_entities.Entity, filter map[string]any) string {

	sql := "SELECT * FROM " + e.GetEntityName() + " WHERE "

	where := ""
	for name, _ := range filter {
		if len(where) > 0 {
			where += " AND "
		}
		where += name + "=?"
	}

	sql += where + ";"

	return sql
}
