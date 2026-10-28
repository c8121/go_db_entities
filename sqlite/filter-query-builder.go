package sqlite

import (
	"log"
	"reflect"

	"github.com/c8121/go_db_entities"
)

type FilterQueryBuilder struct {
}

func (b FilterQueryBuilder) GetFilterStatement(c *go_db_entities.Context, f go_db_entities.Filter) string {

	if len(f.FieldName) == 0 {
		return ""
	}

	sql := f.FieldName

	if go_db_entities.IsSlice(f.Value) {

		if f.Comparison != 0 && f.Comparison != go_db_entities.CompareIn && f.Comparison != go_db_entities.CompareEqual {
			log.Printf("WARN: Ignore invalid comparison for slice-field %s", f.FieldName)
		}

		values := reflect.ValueOf(f.Value)
		sql += " IN("
		for i := 0; i < values.Len(); i++ {
			if i > 0 {
				sql += ","
			}
			sql += "?"
		}
		sql += ")"

	} else {

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
		case go_db_entities.CompareNotEqual:
			sql += " <> "
		default:
			sql += "="
		}

		sql += "?"
	}

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
