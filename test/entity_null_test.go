package test

import (
	"testing"

	"github.com/c8121/go_db_entities"
)

func TestEntityWithNullValues(t *testing.T) {

	ctx := OpenTestDB()
	defer go_db_entities.CloseOrLog(ctx.DB)

	err := go_db_entities.Create(ctx, &Person{})
	if err != nil {
		t.Fatal(err)
	}

	_, err = ctx.DB.Exec("INSERT INTO PERSON (FirstName, LastName) VALUES (NULL, NULL);")
	if err != nil {
		t.Fatal(err)
	}

	itemSupplier := func() go_db_entities.Entity {
		return &Person{}
	}

	filter := go_db_entities.NewFilter("Id", go_db_entities.CompareGreaterEqual, 0)
	_, err = go_db_entities.Select(ctx, itemSupplier, *filter)
	if err != nil {
		t.Fatal(err)
	}

}
