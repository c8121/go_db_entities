package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/c8121/go_db_entities"
	"github.com/c8121/go_db_entities/sqlite"
)

func TestFilter(t *testing.T) {

	ctx := OpenTestDB()
	defer go_db_entities.CloseOrLog(ctx)

	queryBuilder := sqlite.FilterQueryBuilder{}

	filter := go_db_entities.NewFilter("FirstName", go_db_entities.CompareEqual, "Jon")
	query := queryBuilder.GetFilterStatement(ctx, *filter)
	values := filter.GetFilterValues()
	fmt.Printf("Query: %s, Values: %v\n", query, values)
	if !strings.Contains(query, "FirstName") || !strings.Contains(query, "=") || !strings.Contains(query, "?") {
		t.Errorf("Invalid query: %s", query)
	}
	if len(values) != 1 || values[0] != "Jon" {
		t.Errorf("Invalid values: num=%d, %v", len(values), values)
	}

	lastNameFilter := go_db_entities.NewFilter("LastName", go_db_entities.CompareEqual, "Doe")
	filter.Next = lastNameFilter
	query = queryBuilder.GetFilterStatement(ctx, *filter)
	values = filter.GetFilterValues()
	fmt.Printf("Query: %s, Values: %v\n", query, values)
	if !strings.Contains(query, "LastName") || !strings.Contains(query, "AND") {
		t.Errorf("Invalid query: %s", query)
	}
	if len(values) != 2 || values[1] != "Doe" {
		t.Errorf("Invalid values: num=%d, %v", len(values), values)
	}

	lastNameFilter.Next = go_db_entities.NewFilter("LastName", go_db_entities.CompareEqual, "Done")
	lastNameFilter.NextConjunction = go_db_entities.ConjunctOr
	query = queryBuilder.GetFilterStatement(ctx, *filter)
	values = filter.GetFilterValues()
	fmt.Printf("Query: %s, Values: %v\n", query, values)
	if !strings.Contains(query, "OR") || !strings.Contains(query, "?))") {
		t.Errorf("Invalid query: %s", query)
	}
	if len(values) != 3 || values[2] != "Done" {
		t.Errorf("Invalid values: num=%d, %v", len(values), values)
	}
}

func TestFilterWithoutWhere(t *testing.T) {

	ctx := OpenTestDB()
	defer go_db_entities.CloseOrLog(ctx)

	queryBuilder := sqlite.TableQueryBuilder{}

	query := go_db_entities.TableQuery{
		FieldNames: []string{"FirstName", "LastName"},
		Entity:     &Person{},
		Filter: &go_db_entities.Filter{
			Offset: 0,
			Count:  10,
		},
	}

	sql := queryBuilder.GetTableQueryStatement(ctx, query)
	fmt.Printf("Query: %s\n", sql)
	if !strings.Contains(sql, "LIMIT") || strings.Contains(sql, "WHERE") {
		t.Errorf("Invalid query: %s", sql)
	}
}

func TestFilterWithOrderBy(t *testing.T) {

	ctx := OpenTestDB()
	defer go_db_entities.CloseOrLog(ctx)

	err := go_db_entities.Create(ctx, &Person{})
	if err != nil {
		t.Fatal(err)
	}

	err = CreateTestPersons(ctx, 100)
	if err != nil {
		t.Fatal(err)
	}

	queryBuilder := sqlite.TableQueryBuilder{}

	query := go_db_entities.TableQuery{
		FieldNames: []string{"FirstName", "LastName"},
		Entity:     &Person{},
		Filter: &go_db_entities.Filter{
			OrderBy: []string{"FirstName", "LastName"},
		},
	}

	sql := queryBuilder.GetTableQueryStatement(ctx, query)
	fmt.Printf("Query: %s\n", sql)
	if !strings.Contains(sql, "ORDER BY") || strings.Contains(sql, "WHERE") {
		t.Errorf("Invalid query: %s", sql)
	}

	query.Filter.FieldName = "FirstName"
	query.Filter.Value = "Jon_99"
	sql = queryBuilder.GetTableQueryStatement(ctx, query)
	fmt.Printf("Query: %s\n", sql)
	if !strings.Contains(sql, "ORDER BY") || !strings.Contains(sql, "WHERE") {
		t.Errorf("Invalid query: %s", sql)
	}

	query.Filter.Count = 1
	sql = queryBuilder.GetTableQueryStatement(ctx, query)
	fmt.Printf("Query: %s\n", sql)
	if !strings.Contains(sql, "LIMIT") {
		t.Errorf("Invalid query: %s", sql)
	}

	rows, err := go_db_entities.Query(ctx, query)
	if err != nil {
		t.Fatal(err)
	}

	for i, row := range rows {
		fmt.Println(i, row)
	}
}

func TestFilterWithInClause(t *testing.T) {

	ctx := OpenTestDB()
	defer go_db_entities.CloseOrLog(ctx)

	err := go_db_entities.Create(ctx, &Person{})
	if err != nil {
		t.Fatal(err)
	}

	err = CreateTestPersons(ctx, 100)
	if err != nil {
		t.Fatal(err)
	}

	queryBuilder := sqlite.TableQueryBuilder{}

	query := go_db_entities.TableQuery{
		FieldNames: []string{"FirstName", "LastName"},
		Entity:     &Person{},
		Filter: &go_db_entities.Filter{
			FieldName: "FirstName",
			Value:     []string{"Jon_1", "Jon_99"},
			OrderBy:   []string{"FirstName", "LastName"},
		},
	}

	sql := queryBuilder.GetTableQueryStatement(ctx, query)
	fmt.Printf("Query: %s\n", sql)
	if !strings.Contains(sql, "ORDER BY") || !strings.Contains(sql, "WHERE") {
		t.Errorf("Invalid query: %s", sql)
	}

	rows, err := go_db_entities.Query(ctx, query)
	if err != nil {
		t.Fatal(err)
	}

	for i, row := range rows {
		fmt.Println(i, row)
	}
}
