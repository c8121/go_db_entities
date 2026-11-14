package test

import (
	"errors"
	"testing"
	"time"

	"github.com/c8121/go_db_entities"
)

func TestEntityLSD(t *testing.T) {

	storage := OpenTestDB()
	defer go_db_entities.CloseOrLog(storage.DB)

	session := go_db_entities.NewSession(storage)

	person := &Person{}

	err := go_db_entities.Create(session, person)
	if err != nil {
		t.Fatal(err)
	}

	err = go_db_entities.Load(session, person, *go_db_entities.NewIdFilter(1))
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

	err = go_db_entities.Save(session, person)
	if err != nil {
		t.Fatal(err)
	}

	person.FirstName = "Jane"
	err = go_db_entities.Save(session, person)
	if err != nil {
		t.Fatal(err)
	}

	id := person.Id
	person = &Person{}
	err = go_db_entities.Load(session, person, *go_db_entities.NewIdFilter(id))
	if err != nil {
		t.Fatal(err)
	}

	if person.FirstName != "Jane" {
		t.Errorf("FirstName should be Jane")
	}

	err = go_db_entities.Delete(session, person)
	if err != nil {
		t.Fatal(err)
	}

	err = go_db_entities.Load(session, person, *go_db_entities.NewIdFilter(id))
	if !errors.Is(err, go_db_entities.ErrNotFound) {
		t.Errorf("Load() should return ErrNotFound")
	}
}
