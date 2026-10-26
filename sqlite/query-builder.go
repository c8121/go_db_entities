package sqlite

import "github.com/c8121/go_db_entities"

func NewQueryBuilder() *go_db_entities.QueryBuilder {
	return &go_db_entities.QueryBuilder{
		CreateQueryBuilder: &CreateQueryBuilder{},
		InsertQueryBuilder: &InsertQueryBuilder{},
		UpdateQueryBuilder: &UpdateQueryBuilder{},
		DeleteQueryBuilder: &DeleteQueryBuilder{},
		SelectQueryBuilder: &SelectQueryBuilder{},
		FilterQueryBuilder: &FilterQueryBuilder{},
	}
}
