package test

import (
	"testing"

	"github.com/c8121/go_db_entities"
)

func TestEntityWithNullValues(t *testing.T) {

	storage := OpenTestDB()
	defer go_db_entities.CloseOrLog(storage.DB)

	session := go_db_entities.NewSession(storage)

	err := go_db_entities.Create(session, &Person{})
	if err != nil {
		t.Fatal(err)
	}

	stmt, err := storage.DB.Prepare("INSERT INTO PERSON (FirstName, LastName, CreatedAt, UpdatedAt) VALUES (NULL, NULL, NULL, NULL);")
	if err != nil {
		t.Fatal(err)
	}
	defer go_db_entities.CloseOrLog(stmt)

	_, err = stmt.Exec()
	if err != nil {
		t.Fatal(err)
	}

	itemSupplier := func() go_db_entities.Entity {
		return &Person{}
	}

	filter := go_db_entities.NewFilter("Id", go_db_entities.CompareGreaterEqual, 0)
	_, err = go_db_entities.Select(session, itemSupplier, *filter)
	if err != nil {
		t.Fatal(err)
	}

}
