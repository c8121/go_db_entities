package test

import (
	"fmt"
	"testing"

	"github.com/c8121/go_db_entities"
)

func TestEntitySelect(t *testing.T) {

	storage := OpenTestDB()
	defer go_db_entities.CloseOrLog(storage.DB)

	session := go_db_entities.NewSession(storage)

	err := go_db_entities.Create(session, &Person{})
	if err != nil {
		t.Fatal(err)
	}

	err = CreateTestPersons(session, 100)
	if err != nil {
		t.Fatal(err)
	}

	itemSupplier := func() go_db_entities.Entity {
		return &Person{}
	}

	filter := go_db_entities.NewFilter("Id", go_db_entities.CompareGreaterEqual, 10)
	filter.Offset = 10
	filter.Count = 10

	list, err := go_db_entities.Select(session, itemSupplier, *filter)
	if err != nil {
		t.Fatal(err)
	}

	for i, person := range list {
		fmt.Printf("%d: %v\n", i, person)
	}
}
