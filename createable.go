package go_db_entities

import "log"

// CreateQueryBuilder generates DDL-Statements.
type CreateQueryBuilder interface {
	GetCreateQuery(c *Context, e Entity) string
}

// Create takes an Entity and creates the table in the connected database.
// Uses CreateQueryBuilder to generate the DDL-Statement.
func Create(c *Context, e Entity) error {

	if e == nil {
		return ErrNil
	}

	sql := c.QueryBuilder.CreateQueryBuilder.GetCreateQuery(c, e)
	log.Println(sql)

	_, err := c.DB.Exec(sql)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
