package go_db_entities

import (
	"log"
)

const (
	InnerJoin = 1
	LeftJoin  = 2
)

// TableQueryBuilder generates SQL-Select-Statements.
type TableQueryBuilder interface {
	GetTableQueryStatement(c *Context, q TableQuery) string
}

type TableQuery struct {
	FieldNames []string
	Entity     Entity
	Filter     *Filter

	Join []TableQueryJoin
}

type TableQueryJoin struct {
	TableQuery *TableQuery
	Type       int
	OnFrom     string
	OnTo       string
}

func (q *TableQuery) AddJoin(join TableQueryJoin) {

	if q.Join == nil {
		q.Join = make([]TableQueryJoin, 0)
	}
	q.Join = append(q.Join, join)
}

func Query(c *Context, q TableQuery) ([][]any, error) {

	sql := c.QueryBuilder.TableQueryBuilder.GetTableQueryStatement(c, q)
	log.Println(sql)

	stmt, err := c.DB.Prepare(sql)
	if err != nil {
		log.Printf("Error preparing statement: %s", err)
		return nil, err
	}
	defer CloseOrLog(stmt)

	args := q.GetFilterValues()
	log.Println(args)

	rows, err := stmt.Query(args...)
	if err != nil {
		log.Printf("Error executing query: %s", err)
		return nil, err
	}
	defer CloseOrLog(rows)

	columns, err := rows.Columns()
	if err != nil {
		log.Printf("Error getting columns: %s", err)
		return nil, err
	}

	var result = make([][]any, 0)

	for rows.Next() {

		scanArgs := make([]interface{}, len(columns))
		for i := range scanArgs {
			scanArgs[i] = &scanArgs[i]
		}

		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		result = append(result, scanArgs)
	}

	if len(result) == 0 {
		return result, ErrNotFound
	}

	return result, nil
}

func (q TableQuery) GetFilterValues() []any {

	values := make([]any, 0)

	if q.Filter != nil {
		values = append(values, q.Filter.GetFilterValues()...)
	}

	if q.Join != nil {
		for _, join := range q.Join {
			if join.TableQuery.Filter != nil {
				values = append(values, join.TableQuery.Filter.GetFilterValues()...)
			}
		}
	}

	return values
}
