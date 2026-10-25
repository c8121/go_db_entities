package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/c8121/go_db_entities"
)

func TestEntitySelect(t *testing.T) {

	db := OpenTestDB()
	defer go_db_entities.CloseOrLog(db)

	err := go_db_entities.Create(&Person{})
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 100; i++ {
		person := &Person{
			FirstName:    fmt.Sprintf("Jon_%d", i),
			LastName:     fmt.Sprintf("Doe_%d", i),
			StreetNumber: 1,
			CreatedAt:    time.Time{},
			UpdatedAt:    time.Time{},
		}

		err = go_db_entities.Save(person)
		if err != nil {
			t.Fatal(err)
		}
	}

	filter := map[string]any{}
	filter["StreetNumber"] = 1

	itemSupplier := func() go_db_entities.Entity {
		return &Person{}
	}

	list, err := go_db_entities.Select(itemSupplier, filter)
	if err != nil {
		t.Fatal(err)
	}

	for i, person := range list {
		fmt.Printf("%d: %v\n", i, person)
	}
}
