package go_db_entities

import (
	"log"
	"strconv"
	"strings"
)

const (
	StructTagName = "dbe"
)

type FieldProperties struct {
	Name  string
	Size  int
	Index bool
}

// CreateQueryBuilder generates DDL-Statements.
type CreateQueryBuilder interface {
	GetCreateQueries(c *Context, e Entity) []string
}

// Create takes an Entity and creates the table in the connected database.
// Uses CreateQueryBuilder to generate the DDL-Statement.
func Create(c *Context, e Entity) error {

	if e == nil {
		return ErrNil
	}

	ddl := c.QueryBuilder.CreateQueryBuilder.GetCreateQueries(c, e)
	for _, query := range ddl {
		log.Println(query)

		_, err := c.DB.Exec(query)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

// GetProperties reads tags from struct fields to create FieldProperties
func GetProperties(c *Context, e Entity) map[string]FieldProperties {

	properties := make(map[string]FieldProperties)

	entityFields := GetVisibleFields(e)
	for _, field := range entityFields {
		tag := field.Tag.Get(StructTagName)
		if tag != "" {
			props := FieldProperties{
				Name: field.Name,
			}

			items := strings.Split(tag, ",")
			for _, item := range items {
				item = strings.ToLower(strings.TrimSpace(item))
				if item == "index" {
					props.Index = true
				} else if strings.HasPrefix(item, "size:") {
					var err error
					props.Size, err = strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(item, "size:")))
					if err != nil {
						log.Printf("Error in tag \"%s\", cannot parse size: %s", tag, err)
					}
				}
			}

			properties[field.Name] = props
		}
	}

	return properties
}
