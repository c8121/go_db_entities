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
