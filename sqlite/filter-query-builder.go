package sqlite

import (
	"github.com/c8121/go_db_entities"
)

type FilterQueryBuilder struct {
}

func (b FilterQueryBuilder) GetFilterStatement(c *go_db_entities.Context, f go_db_entities.Filter) string {

	sql := f.FieldName

	switch f.Comparison {
	case go_db_entities.CompareEqual:
		sql += "="
	case go_db_entities.CompareGreater:
		sql += ">"
	case go_db_entities.CompareGraterEqual:
		sql += ">="
	case go_db_entities.CompareLess:
		sql += "<"
	case go_db_entities.CompareLessEqual:
		sql += "<="
	case go_db_entities.CompareContains:
		sql += " like "
	default:
		sql += "="
	}

	sql += "?"

	if f.Next != nil {
		switch f.NextConjunction {
		case go_db_entities.ConjunctAnd:
			sql += " AND "
		case go_db_entities.ConjunctOr:
			sql += " OR "
		default:
			sql += " AND "
		}
		sql += b.GetFilterStatement(c, *f.Next)
		sql = "(" + sql + ")"
	}

	return sql
}
