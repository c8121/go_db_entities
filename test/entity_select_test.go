package test

import (
	"fmt"
	"testing"

	"github.com/c8121/go_db_entities"
)

func TestEntitySelect(t *testing.T) {

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

	itemSupplier := func() go_db_entities.Entity {
		return &Person{}
	}

	filter := go_db_entities.NewFilter("Id", go_db_entities.CompareGraterEqual, 10)
	filter.Offset = 10
	filter.Count = 10

	list, err := go_db_entities.Select(ctx, itemSupplier, *filter)
	if err != nil {
		t.Fatal(err)
	}

	for i, person := range list {
		fmt.Printf("%d: %v\n", i, person)
	}
}
