package test

import (
	"fmt"
	"testing"

	"github.com/c8121/go_db_entities"
)

func TestEntityTableQuery(t *testing.T) {

	storage := OpenTestDB()
	defer go_db_entities.CloseOrLog(storage.DB)

	session := go_db_entities.NewSession(storage)

	err := go_db_entities.Create(session, &Person{})
	if err != nil {
		t.Fatal(err)
	}

	err = go_db_entities.Create(session, &PersonContactChannel{})
	if err != nil {
		t.Fatal(err)
	}

	err = go_db_entities.Create(session, &PersonCategory{})
	if err != nil {
		t.Fatal(err)
	}

	err = CreateTestPersons(session, 100)
	if err != nil {
		t.Fatal(err)
	}

	personQuery := &go_db_entities.TableQuery{
		FieldNames: []string{"*"},
		Entity:     &Person{},
		Filter:     go_db_entities.NewFilter("Person.Id", go_db_entities.CompareGreater, 90),
	}
	sql := storage.QueryBuilder.TableQueryBuilder.GetTableQueryStatement(session, *personQuery)
	fmt.Println(sql)

	numberQuery := &go_db_entities.TableQuery{
		FieldNames: []string{"*"},
		Entity:     &PersonContactChannel{},
	}
	personQuery.AddJoin(go_db_entities.TableQueryJoin{
		TableQuery: numberQuery,
		Type:       go_db_entities.LeftJoin,
		OnFrom:     "Id",
		OnTo:       "PersonId",
	})
	sql = storage.QueryBuilder.TableQueryBuilder.GetTableQueryStatement(session, *personQuery)
	fmt.Println(sql)

	categoryQuery := &go_db_entities.TableQuery{
		FieldNames: []string{"Name"},
		Entity:     &PersonCategory{},
	}
	personQuery.AddJoin(go_db_entities.TableQueryJoin{
		TableQuery: categoryQuery,
		Type:       go_db_entities.LeftJoin,
		OnFrom:     "CategoryId",
		OnTo:       "Id",
	})
	sql = storage.QueryBuilder.TableQueryBuilder.GetTableQueryStatement(session, *personQuery)
	fmt.Println(sql)

	personQuery.Filter.Offset = 0
	personQuery.Filter.Count = 3
	sql = storage.QueryBuilder.TableQueryBuilder.GetTableQueryStatement(session, *personQuery)
	fmt.Println(sql)

	rows, err := go_db_entities.Query(session, *personQuery)
	if err != nil {
		t.Fatal(err)
	}

	for i, row := range rows {
		fmt.Println(i, row)
	}
}
