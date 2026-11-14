# go_db_entities

A small, dependency-free (despite sqlite) Go library that maps plain structs to SQL tables.
No code generation, no query DSL, no SQL strings in your application code.

```go
type Person struct {
    Id        int64
    FirstName string
    LastName  string
}

go_db_entities.Create(session, &Person{})                       // CREATE TABLE
go_db_entities.Save(session, &Person{FirstName: "John"})        // INSERT (or UPDATE)
go_db_entities.Load(session, p, *go_db_entities.NewIdFilter(1)) // SELECT ... WHERE Id=?
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
  pooling, drivers and raw queries stay under your control. Transactions are
  plain `database/sql` transactions.
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
- Transactions across several operations (`NewSessionTx`)
- NULL values in the database are read as zero values
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

    storage := go_db_entities.NewStorage(db, sqlite.NewQueryBuilder(), sqlite.NewFieldMapping())
    session := go_db_entities.NewSession(storage)

    // Create the table (CREATE TABLE IF NOT EXISTS)
    if err := go_db_entities.Create(session, &Person{}); err != nil {
        log.Fatal(err)
    }

    // Insert: Id is 0, so a new row is created and the Id is set
    person := &Person{FirstName: "John", LastName: "Doe"}
    if err := go_db_entities.Save(session, person); err != nil {
        log.Fatal(err)
    }

    // Load it again
    loaded := &Person{}
    if err := go_db_entities.Load(session, loaded, *go_db_entities.NewIdFilter(person.Id)); err != nil {
        log.Fatal(err)
    }
    log.Printf("%+v", loaded)
}
```

## Usage

### Storage and Session

Two small types make up the entry point of the library:

- **`Storage`** bundles the database connection (`storage.DB`), the SQL builders
  for your database and the field mappings. Create one at startup and share it.
- **`Session`** says where statements run. All functions of this package take
  a `*Session` as their first parameter.

```go
storage := go_db_entities.NewStorage(db, sqlite.NewQueryBuilder(), sqlite.NewFieldMapping())

// Statements run directly on the database, each one is committed on its own
session := go_db_entities.NewSession(storage)

// Statements run within a transaction, see "Transactions"
session, tx, err := go_db_entities.NewSessionTx(storage)
```

`NewStorage` panics if `db` or the query builder is `nil`. Types without a
registered mapping use the `DefaultFieldMapping`, which passes values through
unchanged.

A `Storage` holds no mutable state of its own. It and the sessions created by
`NewSession` can be shared by goroutines, as long as you do not modify the field
mappings afterwards. A session from `NewSessionTx` is bound to one transaction
and must not be used by several goroutines at the same time.

The `Storage` does not own the connection, so closing `storage.DB` is up to you.

`Session` is a plain struct with the fields `Storage` and `StatementProvider`.
The latter is anything that has a `Prepare(query string) (*sql.Stmt, error)`
method: a `*sql.DB`, a `*sql.Tx` or a `*sql.Conn`. This lets you run the library
on a connection or transaction you began yourself.

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
err := go_db_entities.Create(session, &Person{})

// Insert (Id == 0) or update (Id != 0)
err = go_db_entities.Save(session, person)

// Or call explicitly
err = go_db_entities.Insert(session, person) // ErrExists if Id != 0
err = go_db_entities.Update(session, person) // ErrNotFound if Id == 0 or no row was updated

// Delete one entity (Id is reset to 0 afterwards)
err = go_db_entities.Delete(session, person)

// Delete many
err = go_db_entities.DeleteBatch(session, &Person{},
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
err := go_db_entities.Validate(session, &Person{})

// Only check, change nothing
for _, result := range go_db_entities.CheckFieldsExist(session, &Person{}) {
    if result.Error != nil {
        fmt.Println("missing or unreadable:", result.FieldName, result.Error)
    }
}
```

`CheckFieldsExist` runs a small `SELECT` for every field and reports the name
of the struct field (`FieldName`). Any error counts as a missing column, not
only "no such column". The table itself must already exist (call `Create`
first). `Validate` never renames, drops or changes the
type of a column.

### Load one entity

```go
person := &Person{}
err := go_db_entities.Load(session, person, *go_db_entities.NewIdFilter(1))
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

entities, err := go_db_entities.Select(session, supplier,
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

### NULL values

A column that contains NULL is read as the zero value of the struct field
(`""` for strings, `0` for numbers, `time.Time{}` for times), so rows written
by other programs can be loaded without errors. The other direction does not
exist: `Insert` and `Update` always write the value of the field, so there is
no way to write NULL, and a zero value cannot be told apart from NULL after
reading. `Query` returns raw values, so it gives you `nil` for NULL.

A custom mapping's `ConvertDbValue` is also called for NULL columns, with `nil`
as value. Return `nil, nil` to get the zero value of the field.

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

rows, err := go_db_entities.Query(session, *personQuery)
for _, row := range rows {
    fmt.Println(row...) // LastName, Name, Value
}
```

Order, offset, count and `DISTINCT` are read from the filter of the main
query. Filters of joined queries are added to the `WHERE` clause.

### Transactions

By default every statement runs directly on the database and is committed on
its own. To run several operations in one transaction, create the session with
`NewSessionTx`. It begins a transaction on `storage.DB` and returns the session
together with the `*sql.Tx`. Committing and rolling back is up to you:

```go
session, tx, err := go_db_entities.NewSessionTx(storage)
if err != nil {
    return err
}
defer tx.Rollback() // no effect after Commit (returns sql.ErrTxDone, which can be ignored)

if err := go_db_entities.Save(session, order); err != nil {
    return err // the deferred Rollback undoes everything
}
if err := go_db_entities.Save(session, item); err != nil {
    return err
}

return tx.Commit()
```

If you need more control (isolation level, read-only, a transaction that is
bound to a `context.Context`), begin the transaction yourself and wrap it:

```go
tx, err := storage.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
if err != nil {
    return err
}
defer tx.Rollback()

session := &go_db_entities.Session{Storage: storage, StatementProvider: tx}
// ... go_db_entities.Save(session, ...)

return tx.Commit()
```

Things to keep in mind:

- **Use only the transaction's session inside the transaction.** A session from
  `NewSession` (or `storage.DB` itself) runs on another connection, outside the
  transaction. With SQLite this can lead to "database is locked" errors, or to
  a deadlock if the pool is limited to one connection.
- **A finished transaction cannot be used anymore.** After `Commit` or
  `Rollback`, its session returns `sql.ErrTxDone`.
- **Entities are not rolled back.** `Insert` sets the Id of an entity and
  `Delete` resets it to 0. After a rollback the entities in memory still carry
  these changes: a rolled-back insert leaves the Id set, so a later `Save`
  would try to `UPDATE` a row that does not exist and return `ErrNotFound`.
  Reset the Id or reload the entity.
- **DDL and transactions.** `Create` and `Validate` can run within a
  transaction. SQLite rolls such changes back, but other databases (for example
  MySQL) commit DDL statements implicitly.
- **No nested transactions or savepoints.**

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

Mappings are registered by type name in the map passed to `NewStorage`. The
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
`ErrNotFound` for an empty result set. Errors of `database/sql` (for example
`sql.ErrTxDone`) are returned as they are.

### Logging

The library logs through the standard `log/slog` package. At debug level it
logs every SQL statement and its parameter values. Failures while preparing
or executing statements are logged at error level before the error is
returned. Do not enable debug logging in production if your data is
sensitive.

## How it works

`Storage` holds a `QueryBuilder`, a struct of small interfaces, one per kind
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
lists and runs the statements through the `StatementProvider` of the `Session`
(a `*sql.DB` or a `*sql.Tx`). The SQL text itself comes from the builders,
which receive the `*Session` and only generate SQL; they never access the
database. To support another database, implement these interfaces in a new
package (see `sqlite/` as a reference) and pass the result to `NewStorage`.
The generated SQL currently uses `?` placeholders.

`Validate` first asks the database which columns are missing
(`CheckFieldsExist`), then hands their names to
`CreateQueryBuilder.GetAddFieldsQueries`, which returns the statements to run.

The reflection helpers `GetProperties` (reads the struct tags) and
`GetVisibleFields` are exported, so builders for other databases can use
them as the SQLite builders do.

## Upgrading from 0.0.x

Version 0.1.0 introduced `Storage` and `Session` in place of `Context`, added
transactions and reads NULL values:

| Before                                          | Now                                                        |
|-------------------------------------------------|------------------------------------------------------------|
| `NewContext(db, qb, mappings)`                  | `NewStorage(db, qb, mappings)` and `NewSession(storage)`   |
| type `Context`                                  | `Storage` (database, builders, mappings) and `Session` (where statements run) |
| `ctx.DB`                                        | `storage.DB`                                               |
| `Save(ctx, e)`, `Load(ctx, e, f)`, ... (all functions) | the same calls, with a `*Session` as first parameter |
| builder interfaces take `*Context`              | builder interfaces take `*Session`                         |
| `Context.GetTypeMapping`                        | `Storage.GetTypeMapping`                                   |
| `CreateQueryBuilder.GetCreateMissingFieldsQueries(c, e)` | `GetAddFieldsQueries(session, e, fieldNames)`     |
| `CheckFieldResult.FielName`                     | `CheckFieldResult.FieldName`                               |
| reading NULL panicked                           | NULL is read as the zero value                             |

The `CreateQueryBuilder` row matters only if you implemented your own builder:
it no longer accesses the database; `Validate` does the check and passes the
names of the missing fields.

## Current limitations

This project is intentionally small. Be aware of the following before using it:

- **Supported column types are limited** to integers, `string`, `time.Time`
  and custom mappings. Other types (`bool`, floats, `[]byte`, pointers,
  named types) are not supported yet, and `Create` panics on them.
- **Struct layout:** entities need an exported `int64` field named `Id`.
  Unexported fields and embedded structs are not supported, and there is no
  tag to skip a field.
- **NULL can be read, but not written.** NULL columns are read as zero values
  and `Insert`/`Update` always write the field's value (see
  [NULL values](#null-values)). The included `Embedding` mapping does not handle
  NULL yet and panics on it.
- **No `context.Context` support.** Statements run without a context, so they
  cannot be canceled and have no deadline. A transaction you begin yourself
  with `BeginTx` can be bound to a context (see [Transactions](#transactions)).
- **Transactions are manual and flat.** There are no nested transactions or
  savepoints, and entities are not restored after a rollback.
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
you need migrations, relations, savepoints, cancelable queries, columns where
NULL must differ from the zero value or support for many databases today. Use `go_db_entities` if
you want a small, understandable layer for storing and loading structs with
as few moving parts as possible.

## Running the tests

```sh
go test ./...
```

The tests use an in-memory SQLite database and need cgo
(for `mattn/go-sqlite3`).
