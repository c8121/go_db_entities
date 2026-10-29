package sqlite

import "github.com/c8121/go_db_entities"

func NewFieldMapping() map[string]go_db_entities.FieldMapping {

	mapping := make(map[string]go_db_entities.FieldMapping, 0)
	mapping["Embedding"] = &EmbeddingFieldMapping{}
	return mapping

}
