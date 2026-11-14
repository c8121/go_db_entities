package test

import (
	"errors"
	"testing"
	"time"

	"github.com/c8121/go_db_entities"
)

func TestEntityCRUD(t *testing.T) {

	storage := OpenTestDB()
	defer go_db_entities.CloseOrLog(storage.DB)

	session := go_db_entities.NewSession(storage)

	person := &Person{
		FirstName:    "John",
		LastName:     "Doe",
		StreetName:   "Example Street",
		StreetNumber: 1,
		Zip:          1234,
		City:         "New York",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := go_db_entities.Create(session, person)
	if err != nil {
		t.Fatal(err)
	}

	err = go_db_entities.Insert(session, person)
	if err != nil {
		t.Fatal(err)
	}

	loadedPerson := &Person{}
	loadedPersonSupplier := func() go_db_entities.Entity { return loadedPerson }

	_, err = go_db_entities.Select(session, loadedPersonSupplier, *go_db_entities.NewIdFilter(person.Id))
	if err != nil {
		t.Error(err)
	}

	if person.Id != loadedPerson.Id {
		t.Errorf("person id should be %d, but %d", person.Id, loadedPerson.Id)
	}
	if person.FirstName != loadedPerson.FirstName {
		t.Errorf("person first name should be %s, but %s", person.FirstName, loadedPerson.FirstName)
	}
	if person.CreatedAt.After(loadedPerson.CreatedAt) {
		t.Errorf("person created at should be %v, but %v", person.CreatedAt, loadedPerson.CreatedAt)
	}

	loadedPerson.FirstName = "Jane"
	err = go_db_entities.Update(session, loadedPerson)
	if err != nil {
		t.Error(err)
	}

	_, err = go_db_entities.Select(session, loadedPersonSupplier, *go_db_entities.NewIdFilter(person.Id))
	if err != nil {
		t.Error(err)
	}
	if loadedPerson.FirstName != "Jane" {
		t.Errorf("person first name should be %s, but %s", "Jane", loadedPerson.FirstName)
	}

	err = go_db_entities.Delete(session, loadedPerson)
	if err != nil {
		t.Error(err)
	}

	_, err = go_db_entities.Select(session, loadedPersonSupplier, *go_db_entities.NewIdFilter(person.Id))
	if err != nil && !errors.Is(err, go_db_entities.ErrNotFound) {
		t.Error(err)
	}
	if loadedPerson.Id != 0 {
		t.Errorf("person id should be %d, but %d", 0, loadedPerson.Id)
	}
}
