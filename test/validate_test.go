package test

import (
	"testing"

	"github.com/c8121/go_db_entities"
)

type ErrorPerson struct {
	Id                       int64
	ExtraField               string
	ExtraIndexedField        string `dbe:"Index"`
	ExtraRenamedField        string `dbe:"Name: ExtraRename"`
	ExtraRenamedIndexedField string `dbe:"Name: ExtraRename2, Index"`
}

func (p *ErrorPerson) GetEntityName() string {
	return "Person"
}

func (p *ErrorPerson) GetId() int64 {
	return p.Id
}

func (p *ErrorPerson) SetId(id int64) {
	p.Id = id
}

func TestEntityValidation(t *testing.T) {

	storage := OpenTestDB()
	defer go_db_entities.CloseOrLog(storage.DB)

	session := go_db_entities.NewSession(storage)

	person := &Person{}
	err := go_db_entities.Create(session, person)
	if err != nil {
		t.Fatal(err)
	}

	CheckFieldsExistTest(t, session, person, false)

	err = go_db_entities.Validate(session, person)
	if err != nil {
		t.Fatal(err)
	}

	errPerson := &ErrorPerson{}
	CheckFieldsExistTest(t, session, errPerson, true)

	err = go_db_entities.Validate(session, errPerson)
	if err != nil {
		t.Fatal(err)
	}
}

func CheckFieldsExistTest(t *testing.T, session *go_db_entities.Session, e go_db_entities.Entity, expectError bool) {
	hadError := false
	fields := go_db_entities.CheckFieldsExist(session, e)
	for _, field := range fields {
		t.Log(field)
		if field.Error != nil {
			hadError = true
			if !expectError {
				t.Fatal(field.Error)
			}
		}
	}

	if !hadError && expectError {
		t.Fatal("Expected error but got none")
	}
}
