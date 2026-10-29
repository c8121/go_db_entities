# go_db_entites

Maps structs to database.

Database agnostic (but currently only SQLite), no SQL required on usage.

## Usage

### Create Context

A Context holds a database connection and query builders.

```go
    
url := "file:TestDb?mode=memory&cache=shared"
db, err := sql.Open("sqlite3", url)
if err != nil {
    panic(err)
}

ctx := go_db_entities.NewContext(db, sqlite.NewQueryBuilder())

```

### Create a struct with required field and functions

In order to make a struct storable into a database, 
create struct with field `Id` and functions `GetEntityName`, `GetId` and `SetId`

```go
type Person struct {
    Id           int64
    FirstName    string
    LastName     string
}

// GetEntityName returns the database table name
func (p *Person) GetEntityName() string {
    return "Person"
}

func (p *Person) GetId() int64 {
    return p.Id
}

func (p *Person) SetId(id int64) {
    p.Id = id
}
```

### Create Database table

If desired, the database table can be created automatically

```go
err := go_db_entities.Create(ctx, person)
if err != nil {
    panic(err)
}
```

### Save to database

```go
person := &Person{
    FirstName:    "John",
    LastName:     "Doe",
}

err := go_db_entities.Save(ctx, person)
if err != nil {
    panic(err)
}
```

### Load one entity from database

```go
id := 1
person := &Person{}
err := go_db_entities.Load(ctx, person, *go_db_entities.NewIdFilter(id))
if err != nil {
    t.Fatal(err)
}
```

### Load multiple rows from database

```go
itemSupplier := func() go_db_entities.Entity {
    return &Person{}
}

list, err := go_db_entities.Select(ctx, itemSupplier, *go_db_entities.NewFilter("Id", go_db_entities.CompareGraterEqual, 0))
if err != nil {
    t.Fatal(err)
}
```

### Create queries over multiple tables

```go
personQuery := &go_db_entities.TableQuery{
    FieldNames: []string{"LastName"},
    Entity:     &Person{},
    Filter:     go_db_entities.NewFilter("Person.Id", go_db_entities.CompareGreater, 90),
}
numberQuery := &go_db_entities.TableQuery{
    FieldNames: []string{"Name", "Value"},
    Entity:     &PersonContactChannel{},
}
personQuery.AddJoin(go_db_entities.TableQueryJoin{
    TableQuery: numberQuery,
    Type:       go_db_entities.LeftJoin,
    OnFrom:     "Id",
    OnTo:       "PersonId",
})

rows, err := go_db_entities.Query(ctx, *personQuery)
```