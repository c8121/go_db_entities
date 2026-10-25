package go_db_entities

import "log"

// CreateQueryBuilder generates DDL-Statements.
type CreateQueryBuilder interface {
	GetCreateQuery(e Entity) string
}

// Create takes an Entity and creates the table in the connected database.
// Uses CreateQueryBuilder to generate the DDL-Statement.
func Create(e Entity) error {

	if e == nil {
		return ErrNil
	}

	sql := createBuilder.GetCreateQuery(e)
	log.Println(sql)

	_, err := db.Exec(sql)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
