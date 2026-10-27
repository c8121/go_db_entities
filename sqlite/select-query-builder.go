package sqlite

import (
	"fmt"
	"log"

	"github.com/c8121/go_db_entities"
)

type SelectQueryBuilder struct {
}

func (b SelectQueryBuilder) GetSelectQuery(c *go_db_entities.Context, e go_db_entities.Entity, filter go_db_entities.Filter) string {

	sql := "SELECT * FROM " + e.GetEntityName() + " WHERE "
	sql += c.QueryBuilder.FilterQueryBuilder.GetFilterStatement(c, filter)

	if filter.Offset > 0 || filter.Count > 0 {
		sql += fmt.Sprintf(" LIMIT %d,%d", filter.Offset, filter.Count)
	}

	if filter.Next != nil && filter.Next.Offset > 0 && filter.Next.Count > 0 {
		log.Printf("Warn: Ignoring offest/count of nested filter")
	}

	sql += ";"
	return sql
}
