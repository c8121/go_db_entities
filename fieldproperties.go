package go_db_entities

import (
	"log/slog"
	"strconv"
	"strings"
)

type FieldProperties struct {
	Name  string
	Type  string
	Size  int
	Index bool
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
				item = strings.TrimSpace(item)
				checkName := strings.ToLower(item)
				if checkName == "index" {
					props.Index = true
				} else if strings.HasPrefix(checkName, "name:") {
					props.Name = strings.TrimSpace(item[len("name:"):])
				} else if strings.HasPrefix(checkName, "type:") {
					props.Type = strings.TrimSpace(item[len("type:"):])
				} else if strings.HasPrefix(checkName, "size:") {
					var err error
					props.Size, err = strconv.Atoi(strings.TrimSpace(item[len("size:"):]))
					if err != nil {
						slog.Error("Error in tag \"%s\", cannot parse size: %s", tag, err)
					}
				}
			}

			properties[field.Name] = props
		}
	}

	return properties
}
