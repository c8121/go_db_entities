package test

import (
	"errors"
	"testing"
	"time"

	"github.com/c8121/go_db_entities"
)

func TestEntityLSD(t *testing.T) {

	ctx := OpenTestDB()
	defer go_db_entities.CloseOrLog(ctx)

	person := &Person{}

	err := go_db_entities.Create(ctx, person)
	if err != nil {
		t.Fatal(err)
	}

	err = go_db_entities.Load(ctx, person, *go_db_entities.NewIdFilter(1))
	if !errors.Is(err, go_db_entities.ErrNotFound) {
		t.Errorf("Load() should return ErrNotFound")
	}

	person = &Person{
		FirstName:    "John",
		LastName:     "Doe",
		StreetName:   "Example Street",
		StreetNumber: 1,
		Zip:          1234,
		City:         "New York",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = go_db_entities.Save(ctx, person)
	if err != nil {
		t.Fatal(err)
	}

	person.FirstName = "Jane"
	err = go_db_entities.Save(ctx, person)
	if err != nil {
		t.Fatal(err)
	}

	id := person.Id
	person = &Person{}
	err = go_db_entities.Load(ctx, person, *go_db_entities.NewIdFilter(id))
	if err != nil {
		t.Fatal(err)
	}

	if person.FirstName != "Jane" {
		t.Errorf("FirstName should be Jane")
	}

	err = go_db_entities.Delete(ctx, person)
	if err != nil {
		t.Fatal(err)
	}

	err = go_db_entities.Load(ctx, person, *go_db_entities.NewIdFilter(id))
	if !errors.Is(err, go_db_entities.ErrNotFound) {
		t.Errorf("Load() should return ErrNotFound")
	}
}
