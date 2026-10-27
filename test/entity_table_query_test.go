package test

import (
	"fmt"
	"testing"

	"github.com/c8121/go_db_entities"
)

func TestEntityTableQuery(t *testing.T) {

	ctx := OpenTestDB()
	defer go_db_entities.CloseOrLog(ctx)

	err := go_db_entities.Create(ctx, &Person{})
	if err != nil {
		t.Fatal(err)
	}

	err = go_db_entities.Create(ctx, &PersonContactChannel{})
	if err != nil {
		t.Fatal(err)
	}

	err = go_db_entities.Create(ctx, &PersonCategory{})
	if err != nil {
		t.Fatal(err)
	}

	err = CreateTestPersons(ctx, 100)
	if err != nil {
		t.Fatal(err)
	}

	personQuery := &go_db_entities.TableQuery{
		FieldNames: []string{"*"},
		Entity:     &Person{},
		Filter:     go_db_entities.NewFilter("Person.Id", go_db_entities.CompareGreater, 90),
	}
	sql := ctx.QueryBuilder.TableQueryBuilder.GetTableQueryStatement(ctx, *personQuery)
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
	sql = ctx.QueryBuilder.TableQueryBuilder.GetTableQueryStatement(ctx, *personQuery)
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
	sql = ctx.QueryBuilder.TableQueryBuilder.GetTableQueryStatement(ctx, *personQuery)
	fmt.Println(sql)

	rows, err := go_db_entities.Query(ctx, *personQuery)
	if err != nil {
		t.Fatal(err)
	}

	for i, row := range rows {
		fmt.Println(i, row)
	}
}
