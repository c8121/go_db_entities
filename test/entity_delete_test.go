package test

import (
	"errors"
	"testing"

	"github.com/c8121/go_db_entities"
)

func TestEntityDelete(t *testing.T) {

	ctx := OpenTestDB()
	defer go_db_entities.CloseOrLog(ctx)

	err := go_db_entities.Create(ctx, &Person{})
	if err != nil {
		t.Fatal(err)
	}

	err = CreateTestPersons(ctx, 3)
	if err != nil {
		t.Fatal(err)
	}

	person := &Person{}
	err = go_db_entities.Load(ctx, person, *go_db_entities.NewIdFilter(1))
	if err != nil {
		t.Fatal(err)
	}

	err = go_db_entities.Delete(ctx, person)
	if err != nil {
		t.Fatal(err)
	}

	items, err := go_db_entities.Select(ctx, func() go_db_entities.Entity { return &Person{} }, go_db_entities.Filter{})
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 2 {
		t.Errorf("Expected 2 items after delete, have %d", len(items))
	}
}

func TestEntityDeleteBatch(t *testing.T) {

	ctx := OpenTestDB()
	defer go_db_entities.CloseOrLog(ctx)

	err := go_db_entities.Create(ctx, &Person{})
	if err != nil {
		t.Fatal(err)
	}

	err = CreateTestPersons(ctx, 3)
	if err != nil {
		t.Fatal(err)
	}

	// Delete one
	person := &Person{}
	err = go_db_entities.DeleteBatch(ctx, person, *go_db_entities.NewIdFilter(1))
	if err != nil {
		t.Fatal(err)
	}

	items, err := go_db_entities.Select(ctx, func() go_db_entities.Entity { return &Person{} }, go_db_entities.Filter{})
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 2 {
		t.Errorf("Expected 99 items after delete, have %d", len(items))
	}

	// Delete all
	err = go_db_entities.DeleteBatch(ctx, person, go_db_entities.Filter{FieldName: "Id", Comparison: go_db_entities.CompareNotEqual, Value: 1})
	if err != nil {
		t.Fatal(err)
	}

	items, err = go_db_entities.Select(ctx, func() go_db_entities.Entity { return &Person{} }, go_db_entities.Filter{})
	if err != nil && !errors.Is(go_db_entities.ErrNotFound, err) {
		t.Fatal(err)
	}

	if len(items) != 0 {
		t.Errorf("Expected 0 items after delete, have %d", len(items))
	}
}
