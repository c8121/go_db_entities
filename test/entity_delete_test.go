package test

import (
	"errors"
	"testing"

	"github.com/c8121/go_db_entities"
)

func TestEntityDelete(t *testing.T) {

	storage := OpenTestDB()
	defer go_db_entities.CloseOrLog(storage.DB)

	session := go_db_entities.NewSession(storage)

	err := go_db_entities.Create(session, &Person{})
	if err != nil {
		t.Fatal(err)
	}

	err = CreateTestPersons(session, 3)
	if err != nil {
		t.Fatal(err)
	}

	person := &Person{}
	err = go_db_entities.Load(session, person, *go_db_entities.NewIdFilter(1))
	if err != nil {
		t.Fatal(err)
	}

	err = go_db_entities.Delete(session, person)
	if err != nil {
		t.Fatal(err)
	}

	items, err := go_db_entities.Select(session, func() go_db_entities.Entity { return &Person{} }, go_db_entities.Filter{})
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 2 {
		t.Errorf("Expected 2 items after delete, have %d", len(items))
	}
}

func TestEntityDeleteBatch(t *testing.T) {

	ctx := OpenTestDB()
	defer go_db_entities.CloseOrLog(ctx.DB)

	session := go_db_entities.NewSession(ctx)

	err := go_db_entities.Create(session, &Person{})
	if err != nil {
		t.Fatal(err)
	}

	err = CreateTestPersons(session, 3)
	if err != nil {
		t.Fatal(err)
	}

	// Delete one
	person := &Person{}
	err = go_db_entities.DeleteBatch(session, person, *go_db_entities.NewIdFilter(1))
	if err != nil {
		t.Fatal(err)
	}

	items, err := go_db_entities.Select(session, func() go_db_entities.Entity { return &Person{} }, go_db_entities.Filter{})
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 2 {
		t.Errorf("Expected 99 items after delete, have %d", len(items))
	}

	// Delete all
	err = go_db_entities.DeleteBatch(session, person, go_db_entities.Filter{FieldName: "Id", Comparison: go_db_entities.CompareNotEqual, Value: 1})
	if err != nil {
		t.Fatal(err)
	}

	items, err = go_db_entities.Select(session, func() go_db_entities.Entity { return &Person{} }, go_db_entities.Filter{})
	if err != nil && !errors.Is(go_db_entities.ErrNotFound, err) {
		t.Fatal(err)
	}

	if len(items) != 0 {
		t.Errorf("Expected 0 items after delete, have %d", len(items))
	}
}
