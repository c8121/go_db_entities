package go_db_entities

import "log"

type CreateQueryBuilder interface {
	GetCreateQuery(e Entity) string
}

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
