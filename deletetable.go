package go_db_entities

import (
	"log"
	"reflect"
)

// DeleteQueryBuilder generates SQL-Statements.
type DeleteQueryBuilder interface {
	GetDeleteQuery(c *Context, e Entity) string
	GetDeleteBatchQuery(c *Context, e Entity, f Filter) string
}

// Delete takes an Entity and deletes it from the connected database.
func Delete(c *Context, e Entity) error {

	if e == nil {
		return ErrNil
	}

	if e.GetId() == 0 {
		return ErrNotFound
	}

	sql := c.QueryBuilder.DeleteQueryBuilder.GetDeleteQuery(c, e)
	log.Println(sql)

	stmt, err := c.DB.Prepare(sql)
	if err != nil {
		log.Printf("Error preparing statement: %s", err)
		return err
	}
	defer CloseOrLog(stmt)

	v := reflect.ValueOf(GetAsStruct(e))

	values := make([]any, 1)
	values[0] = v.FieldByName("Id").Interface()
	log.Println(values)

	_, err = stmt.Exec(values...)
	if err != nil {
		log.Printf("Error executing query: %s", err)
		return err
	}

	e.SetId(0)

	return nil
}

// DeleteBatch takes an Entity and a Filter to delete multiple items from the connected database.
func DeleteBatch(c *Context, e Entity, f Filter) error {

	if e == nil {
		return ErrNil
	}

	sql := c.QueryBuilder.DeleteQueryBuilder.GetDeleteBatchQuery(c, e, f)
	log.Println(sql)

	stmt, err := c.DB.Prepare(sql)
	if err != nil {
		log.Printf("Error preparing statement: %s", err)
		return err
	}
	defer CloseOrLog(stmt)

	args := f.GetFilterValues()
	log.Println(args)

	_, err = stmt.Exec(args...)
	if err != nil {
		log.Printf("Error executing query: %s", err)
		return err
	}

	return nil
}
