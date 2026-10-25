package test

import "time"

type Person struct {
	Id           int64
	FirstName    string
	LastName     string
	StreetName   string
	StreetNumber int //makes no sense as int in real life, just for testing
	Zip          int
	City         string

	CreatedAt time.Time
	UpdatedAt time.Time
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
