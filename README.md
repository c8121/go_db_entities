# go_db_entities

A small, dependency-free (despite sqlite) Go library that maps plain structs to SQL tables.
No code generation, no query DSL, no SQL strings in your application code.

```go
type Person struct {
    Id        int64
    FirstName string
    LastName  string
}

go_db_entities.Create(ctx, &Person{})                       // CREATE TABLE
go_db_entities.Save(ctx, &Person{FirstName: "John"})        // INSERT (or UPDATE)
go_db_entities.Load(ctx, p, *go_db_entities.NewIdFilter(1)) // SELECT ... WHERE Id=?
```

## Why this library?

Most Go ORMs are large frameworks: they bring code generators, their own
query languages, migration tools, relation systems, hooks and a long list of
transitive dependencies. That is the right tool for big applications, but
overkill if all you want is to store and load structs.

`go_db_entities` takes the opposite approach:

- **Lightweight.** A small codebase you can read in an afternoon. When
  something behaves unexpectedly, you can find out why by reading the code.
- **No dependencies.** The library itself uses only the Go standard library
  (`database/sql`, `reflect`, `log/slog`). You choose and import the database
  driver yourself. The only third-party module in `go.mod`
  (`mattn/go-sqlite3`) is used by the tests.
- **Plain structs.** An entity is a struct with an `Id` field and three tiny
  methods. No base type to embed, no registration step, no generated code.
- **Transparent SQL.** The SQL it generates is simple, readable and logged at
  debug level. Values are always passed as bound parameters.
- **Built on `database/sql`.** You pass in your own `*sql.DB`, so connection
  pooling, timeouts, drivers and raw queries stay under your control.
- **Pluggable by design.** SQL generation and type conversion are behind small
  interfaces. SQLite is included; other databases can be added without
  touching the core.
- **Handles non-standard column types.** Custom field mappings convert Go
  values to and from database representations. This is used, for example, to
  store `[]float32` embeddings with the [sqlite-vec](https://github.com/asg017/sqlite-vec)
  extension.

It is **not** a full-featured ORM. See [Current limitations](#current-limitations)
and [When not to use it](#when-not-to-use-it).

## Features

- Create tables and indexes from struct definitions
- Add columns that are missing in an existing table (`Validate`)
- Insert, update, save (insert-or-update), load, select and delete
- Batch delete by filter
- Composable filters: comparisons, `IN` lists, `AND`/`OR` chaining, ordering,
  offset/count, `DISTINCT`
- Queries across multiple tables with `INNER`/`LEFT` joins
- Struct tags for column name, type, size and indexes
- Custom type mappings (included: `Embedding` for sqlite-vec)

## Installation

```sh
go get github.com/c8121/go_db_entities
```

Requires Go 1.25.1 or newer (see `go.mod`). You also need a `database/sql`
driver for your database. The tests use
[`mattn/go-sqlite3`](https://github.com/mattn/go-sqlite3) (requires cgo).

## Quick start

```go
package main

import (
    "database/sql"
    "log"

    _ "github.com/mattn/go-sqlite3"

    "github.com/c8121/go_db_entities"
    "github.com/c8121/go_db_entities/sqlite"
)

type Person struct {
    Id        int64
    FirstName string
    LastName  string
}

func (p *Person) GetEntityName() string { return "Person" } // table name
func (p *Person) GetId() int64          { return p.Id }
func (p *Person) SetId(id int64)        { p.Id = id }

func main() {
    db, err := sql.Open("sqlite3", "file:app.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    ctx := go_db_entities.NewContext(db, sqlite.NewQueryBuilder(), sqlite.NewFieldMapping())

    // Create the table (CREATE TABLE IF NOT EXISTS)
    if err := go_db_entities.Create(ctx, &Person{}); err != nil {
        log.Fatal(err)
    }

    // Insert: Id is 0, so a new row is created and the Id is set
    person := &Person{FirstName: "John", LastName: "Doe"}
    if err := go_db_entities.Save(ctx, person); err != nil {
        log.Fatal(err)
    }

    // Load it again
    loaded := &Person{}
    if err := go_db_entities.Load(ctx, loaded, *go_db_entities.NewIdFilter(person.Id)); err != nil {
        log.Fatal(err)
    }
    log.Printf("%+v", loaded)
}
```

## Usage

### The Context

A `Context` bundles the database connection, the SQL builders for your
database and the field mappings:

```go
ctx := go_db_entities.NewContext(db, sqlite.NewQueryBuilder(), sqlite.NewFieldMapping())
```

`NewContext` panics if `db` or the query builder is `nil`. Types without a
registered mapping use the `DefaultFieldMapping`, which passes values through
unchanged.

A `Context` holds no mutable state of its own and can be shared by goroutines,
as long as you do not modify its mappings afterwards. Note that it is not a
`context.Context`; see [Current limitations](#current-limitations).

### Defining entities

An entity is a struct with an `int64` field named `Id` and three methods:

```go
type Person struct {
    Id        int64
    FirstName string
    LastName  string
}

func (p *Person) GetEntityName() string { return "Person" } // table name
func (p *Person) GetId() int64          { return p.Id }
func (p *Person) SetId(id int64)        { p.Id = id }
```

Every exported field becomes a column. `Id` becomes `integer PRIMARY KEY`.

Supported field types: `int`, `int8`–`int64`, `string`, `time.Time` and, with
a custom mapping, anything you map yourself (see [Custom type mappings](#custom-type-mappings)).

### Struct tags

Column options are set with the `dbe` tag, separated by commas. Option names
are case-insensitive (`Index`, `index` and `INDEX` are the same), and spaces
around values are ignored:

| Option        | Meaning                                         | Example               |
|---------------|-------------------------------------------------|-----------------------|
| `name: <col>` | Use a different column name                     | `dbe:"name: surname"` |
| `type: <T>`   | Override the type used for mapping and DDL      | `dbe:"type: Embedding"` |
| `size: <n>`   | Size for the column type (e.g. vector length)   | `dbe:"size: 512"`     |
| `index`       | Create an index on the column                   | `dbe:"index"`         |

```go
type Document struct {
    Id        int64
    Title     string    `dbe:"index"`
    Embedding []float32 `dbe:"type: Embedding, size: 512"`
}
```

Several options can be combined, for example `dbe:"Index, Size: 100"`. The
value of `type:` is matched case-sensitively against the registered field
mappings. A `size:` that is not a number is logged as an error and ignored.

The tag affects `CREATE`, `INSERT`, `UPDATE` and `SELECT` statements. Names
you write yourself, such as filter fields and `TableQuery` field names, are
used as column names as they are, so use the name from the `name:` option
there if you set one.

Indexes are created as `IDX_<table>_<column>`. Note that SQLite does not
enforce sizes such as `Text(100)`; they are informational.

### Create, save, update, delete

```go
// Create table and indexes
err := go_db_entities.Create(ctx, &Person{})

// Insert (Id == 0) or update (Id != 0)
err = go_db_entities.Save(ctx, person)

// Or call explicitly
err = go_db_entities.Insert(ctx, person) // ErrExists if Id != 0
err = go_db_entities.Update(ctx, person) // ErrNotFound if Id == 0 or no row was updated

// Delete one entity (Id is reset to 0 afterwards)
err = go_db_entities.Delete(ctx, person)

// Delete many
err = go_db_entities.DeleteBatch(ctx, &Person{},
    *go_db_entities.NewFilter("LastName", go_db_entities.CompareEqual, "Doe"))
```

`Delete` does not report an error if no row with that Id exists, and
`DeleteBatch` does not report how many rows were deleted. `DeleteBatch`
needs a filter with a `FieldName`; an empty `Filter{}` produces invalid SQL
and returns a database error.

### Adding missing columns

`Create` only issues `CREATE TABLE IF NOT EXISTS`, so it does not touch a
table that already exists. If you add a field to a struct later, use
`Validate` to add the missing columns (and indexes) to the existing table:

```go
// ALTER TABLE ... ADD COLUMN for every field that has no column yet,
// plus CREATE INDEX IF NOT EXISTS for indexed fields
err := go_db_entities.Validate(ctx, &Person{})

// Only check, change nothing
for _, result := range go_db_entities.CheckFieldsExist(ctx, &Person{}) {
    if result.Error != nil {
        fmt.Println("missing or unreadable:", result.FielName, result.Error)
    }
}
```

`CheckFieldsExist` runs a small `SELECT` for every field. Any error counts as
a missing column, not only "no such column". The table itself must already
exist (call `Create` first). `Validate` never renames, drops or changes the
type of a column.

### Load one entity

```go
person := &Person{}
err := go_db_entities.Load(ctx, person, *go_db_entities.NewIdFilter(1))
if errors.Is(err, go_db_entities.ErrNotFound) {
    // no such row
}
```

`Load` is meant for filters that match exactly one row. If several rows
match, a warning is logged and the entity ends up with the values of the last
one.

### Select many entities

`Select` needs a supplier function that creates a fresh entity for each row:

```go
supplier := func() go_db_entities.Entity { return &Person{} }

entities, err := go_db_entities.Select(ctx, supplier,
    *go_db_entities.NewFilter("Id", go_db_entities.CompareGreaterEqual, 10))
if errors.Is(err, go_db_entities.ErrNotFound) {
    // no rows matched (entities is empty)
} else if err != nil {
    return err
}

for _, e := range entities {
    person := e.(*Person)
    fmt.Println(person.FirstName, person.LastName)
}
```

### Filters

A `Filter` describes one condition and can be chained to the next one:

```go
// (FirstName = 'John' AND (LastName = 'Doe' OR LastName = 'Smith'))
filter := go_db_entities.NewFilter("FirstName", go_db_entities.CompareEqual, "John")

lastName := go_db_entities.NewFilter("LastName", go_db_entities.CompareEqual, "Doe")
lastName.Next = go_db_entities.NewFilter("LastName", go_db_entities.CompareEqual, "Smith")
lastName.NextConjunction = go_db_entities.ConjunctOr

filter.Next = lastName // NextConjunction defaults to AND
```

Functions take a `Filter` by value, so dereference the pointer returned by
`NewFilter` and `NewIdFilter`: `*filter`.

Comparisons:

| Constant               | SQL                                    |
|------------------------|----------------------------------------|
| `CompareEqual`         | `=` (default)                          |
| `CompareNotEqual`      | `<>`                                   |
| `CompareGreater`       | `>`                                    |
| `CompareGreaterEqual`  | `>=`                                   |
| `CompareLess`          | `<`                                    |
| `CompareLessEqual`     | `<=`                                   |
| `CompareContains`      | `LIKE` (add the `%` wildcards yourself)|
| `CompareIn`            | `IN (...)`, used automatically when the value is a slice |

```go
// Id IN (1, 2, 3)
filter := go_db_entities.NewFilter("Id", go_db_entities.CompareIn, []int64{1, 2, 3})
```

If the value is a slice, an `IN` list is always generated; any other
comparison given together with a slice is ignored with a warning.

Ordering, paging and `DISTINCT` are set on the (first) filter:

```go
filter := go_db_entities.Filter{
    FieldName:  "Id",
    Comparison: go_db_entities.CompareGreater,
    Value:      100,
    OrderBy:    []string{"LastName", "FirstName DESC"},
    Offset:     20,
    Count:      10, // always set Count together with Offset
    Distinct:   false,
}
```

Offset and count of chained filters are ignored (with a warning). An empty
`Filter{}` (no `FieldName`) selects all rows.

### Queries over multiple tables

`Query` joins tables and returns raw rows as `[][]any`, in the order of the
selected columns. Field names are column names as they exist in the database;
no tags or type mappings are applied to them. Qualify field names in filters
with the table name (`Person.Id`) to avoid ambiguity.

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
    Type:       go_db_entities.LeftJoin, // or InnerJoin
    OnFrom:     "Id",
    OnTo:       "PersonId",
})

rows, err := go_db_entities.Query(ctx, *personQuery)
for _, row := range rows {
    fmt.Println(row...) // LastName, Name, Value
}
```

Order, offset, count and `DISTINCT` are read from the filter of the main
query. Filters of joined queries are added to the `WHERE` clause.

### Custom type mappings

A `FieldMapping` controls how a field is named in `SELECT` statements, which
SQL parameter is used for it and how values are converted in both directions:

```go
type FieldMapping interface {
    GetFieldName(name string) string
    GetFieldParameter() string
    ConvertStructValue(value any) (any, error)
    ConvertDbValue(value any) (any, error)
}
```

Mappings are registered by type name in the map passed to `NewContext`. The
type name is the Go type name of the field (for example `string`), or the
value of the `type:` tag if there is one.
`sqlite.NewFieldMapping()` registers `Embedding`, which stores a `[]float32`
as a sqlite-vec vector (`vec_f32(?)` on write, `vec_to_json(...)` on read).
This needs the [sqlite-vec](https://github.com/asg017/sqlite-vec) extension
to be loaded in your SQLite connection, and the field needs a `size:` (the
vector length) so that the column can be created as `float[<size>]`.

A new type also needs support in the `CreateQueryBuilder` of your database,
because that is where Go types are translated to column types.

### Errors

| Error                      | Meaning                                              |
|----------------------------|------------------------------------------------------|
| `ErrNil`                   | A `nil` entity was passed                            |
| `ErrNotFound`              | No row found (`Load`, `Select`, `Query`), `Update` without an Id or without an affected row, `Delete` with `Id == 0` |
| `ErrExists`                | `Insert` was called with an entity that already has an Id |
| `ErrMustBePointerToStruct` | The entity supplier did not return a pointer to a struct |

Use `errors.Is` to check for them. Note that `Select` and `Query` return
`ErrNotFound` for an empty result set.

### Logging

The library logs through the standard `log/slog` package. At debug level it
logs every SQL statement and its parameter values. Failures while preparing
or executing statements are logged at error level before the error is
returned. Do not enable debug logging in production if your data is
sensitive.

## How it works

`Context` holds a `QueryBuilder`, a struct of small interfaces, one per kind
of statement:

```
CreateQueryBuilder   DDL for tables and indexes, and for missing columns
InsertQueryBuilder   INSERT
UpdateQueryBuilder   UPDATE
DeleteQueryBuilder   DELETE (single and batch)
SelectQueryBuilder   SELECT for entities
FilterQueryBuilder   WHERE clauses from filters
TableQueryBuilder    SELECT with joins
```

The core package reads struct fields via reflection, builds the argument
lists and runs the statements through `database/sql`. The SQL text itself
comes from the builders. To support another database, implement these
interfaces in a new package (see `sqlite/` as a reference) and pass the
result to `NewContext`. The generated SQL currently uses `?` placeholders.

The reflection helpers `GetProperties` (reads the struct tags) and
`GetVisibleFields` are exported, so builders for other databases can use
them as the SQLite builders do.

## Current limitations

This project is intentionally small. Be aware of the following before using it:

- **Supported column types are limited** to integers, `string`, `time.Time`
  and custom mappings. Other types (`bool`, floats, `[]byte`, pointers,
  named types) are not supported yet, and `Create` panics on them.
- **No NULL support.** Reading a NULL column into a struct field is not
  handled and panics.
- **Struct layout:** entities need an exported `int64` field named `Id`.
  Unexported fields and embedded structs are not supported, and there is no
  tag to skip a field.
- **No transactions** and no `context.Context` support: all operations run
  directly on the `*sql.DB`.
- **Migrations are minimal.** `Create` only issues `CREATE TABLE IF NOT EXISTS`
  and `Validate` only adds missing columns and indexes. Renaming, dropping
  and changing columns is up to you.
- **`Validate` checks struct field names**, not the names set with `name:`.
  For an existing table with renamed fields it would try to add columns
  that already exist and fail, so do not use it for such entities yet.
- **No relations, hooks or lazy loading.**
- **Paging:** always set `Count` when using `Offset`, otherwise no rows are
  returned.
- **Joins:** filters on joined tables are applied in the `WHERE` clause, so a
  `LEFT JOIN` filtered on a joined column behaves like an inner join. Filters
  of nested joins are not applied.
- Ids are `INTEGER PRIMARY KEY` without `AUTOINCREMENT`, so SQLite may reuse
  the Id of a deleted last row.

## Security notes

Values are always passed as bound parameters and are safe against SQL
injection. **Identifiers are not.** Table names, field names, `OrderBy`
entries, join fields and the values of the `name:` and `type:` tags are
inserted into the SQL text as they are. Never pass untrusted input (for
example a sort column from a URL parameter) to these without validating it
against a fixed list of allowed names first.

## When not to use it

Choose a full ORM or a query builder such as GORM, ent, sqlc or Squirrel if
you need migrations, relations, transactions across several operations,
nullable columns or support for many databases today. Use `go_db_entities` if
you want a small, understandable layer for storing and loading structs with
as few moving parts as possible.

## Running the tests

```sh
go test ./...
```

The tests use an in-memory SQLite database and need cgo
(for `mattn/go-sqlite3`).
