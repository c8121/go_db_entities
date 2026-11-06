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

type EntityProperties struct {
	fieldProperties map[string]FieldProperties
}

func (p *EntityProperties) GetFieldProperties(name string) (FieldProperties, bool) {
	properties, ok := p.fieldProperties[name]
	return properties, ok
}

func (p *EntityProperties) GetFieldName(name string) string {
	properties, ok := p.fieldProperties[name]
	if ok && properties.Name != "" {
		return properties.Name
	}
	return name
}

func (p *EntityProperties) GetFieldType(name string, defaultName string) string {
	properties, ok := p.fieldProperties[name]
	if ok && properties.Type != "" {
		return properties.Type
	}
	return defaultName
}

func (p *EntityProperties) GetFieldSize(name string) int {
	properties, ok := p.fieldProperties[name]
	if ok {
		return properties.Size
	}
	return 0
}

func (p *EntityProperties) GetFieldIndex(name string) bool {
	properties, ok := p.fieldProperties[name]
	return ok && properties.Index
}

// GetProperties reads tags from struct fields to create FieldProperties
func GetProperties(c *Context, e Entity) *EntityProperties {

	properties := &EntityProperties{
		fieldProperties: make(map[string]FieldProperties),
	}

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
						slog.Error("Error in tag, cannot parse size", "tag", tag, "error", err)
					}
				}
			}

			properties.fieldProperties[field.Name] = props
		}
	}

	return properties
}
