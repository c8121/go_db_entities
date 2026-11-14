package test

import (
	"testing"

	"github.com/c8121/go_db_entities"
)

func TestEntityCRUDTransaction(t *testing.T) {

	storage := OpenTestDB()
	defer go_db_entities.CloseOrLog(storage.DB)

	session, tx, err := go_db_entities.NewSessionTx(storage)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := tx.Commit()
		if err != nil {
			t.Fatal(err)
		}
	}()

	person := &Person{
		FirstName: "John",
		LastName:  "Doe",
	}

	err = go_db_entities.Create(session, person)
	if err != nil {
		t.Fatal(err)
	}

	err = go_db_entities.Insert(session, person)
	if err != nil {
		t.Fatal(err)
	}

	err = go_db_entities.Load(session, person, *go_db_entities.NewFilter("Id", go_db_entities.CompareGreaterEqual, 0))
	if err != nil {
		t.Fatal(err)
	}

	err = go_db_entities.Delete(session, person)
	if err != nil {
		t.Fatal(err)
	}
}
