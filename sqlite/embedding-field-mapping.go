package sqlite

import (
	"encoding/json"
)

type EmbeddingFieldMapping struct {
}

func (m *EmbeddingFieldMapping) GetFieldParameter() string {
	return "vec_f32(?)"
}

func (m *EmbeddingFieldMapping) ConvertStructValue(value any) (any, error) {
	embeddingJson, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return string(embeddingJson), nil
}

func (m *EmbeddingFieldMapping) ConvertDbValue(value any) (any, error) {
	var embeddingJson []float32
	if err := json.Unmarshal([]byte(value.(string)), &embeddingJson); err != nil {
		return nil, err
	}
	return embeddingJson, nil
}
