package test

import (
	"fmt"
	"time"

	"github.com/c8121/go_db_entities"
)

type Person struct {
	Id           int64
	FirstName    string `dbe:"Index,Size:100"`
	LastName     string `dbe:"Index"`
	StreetName   string `dbe:"Size: 125"`
	StreetNumber int    //makes no sense as int in real life, just for testing
	Zip          int    `dbe:"Index"`
	City         string

	CreatedAt time.Time
	UpdatedAt time.Time

	CategoryId int64 `dbe:"Index"`
}

func (p *Person) GetEntityName() string {
	return "Person"
}

func (p *Person) GetId() int64 {
	return p.Id
}

func (p *Person) SetId(id int64) {
	p.Id = id
}

type PersonContactChannel struct {
	Id       int64
	PersonId int64  `dbe:"Index"`
	Name     string `dbe:"Size:100"`
	Value    string
}

func (p *PersonContactChannel) GetEntityName() string {
	return "PersonContactChannel"
}

func (p *PersonContactChannel) GetId() int64 {
	return p.Id
}

func (p *PersonContactChannel) SetId(id int64) {
	p.Id = id
}

type PersonCategory struct {
	Id   int64
	Name string `dbe:"Size:100"`
}

func (p *PersonCategory) GetEntityName() string {
	return "PersonCategory"
}

func (p *PersonCategory) GetId() int64 {
	return p.Id
}

func (p *PersonCategory) SetId(id int64) {
	p.Id = id
}

func CreateTestPersons(ctx *go_db_entities.Context, num int) error {

	for i := 0; i < num; i++ {
		person := &Person{
			FirstName:    fmt.Sprintf("Jon_%d", i),
			LastName:     fmt.Sprintf("Doe_%d", i),
			StreetNumber: 1,
			CreatedAt:    time.Time{},
			UpdatedAt:    time.Time{},
		}

		err := go_db_entities.Save(ctx, person)
		if err != nil {
			return err
		}
	}

	return nil
}
