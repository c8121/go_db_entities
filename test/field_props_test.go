package test

import (
	"strings"
	"testing"

	"github.com/c8121/go_db_entities"
)

type TestEntity struct {
	Id int64

	Renamed string    `dbe:"name: Name"`
	Retyped []float32 `dbe:"type: Embedding,Size: 512"`
}

func (t *TestEntity) GetEntityName() string {
	return "TestEntity"
}

func (t *TestEntity) GetId() int64 {
	return t.Id
}

func (t *TestEntity) SetId(id int64) {
	t.Id = id
}

func TestFieldProps(t *testing.T) {

	ctx := OpenTestDB()
	defer go_db_entities.CloseOrLog(ctx)

	sql := ctx.QueryBuilder.CreateQueryBuilder.GetCreateQueries(ctx, &TestEntity{})[0]
	if !strings.Contains(sql, " Name ") {
		t.Errorf("TestFieldProps failed: 'Renamed' not changed to 'Name'")
	}
	if strings.Contains(sql, "Embedding") {
		t.Errorf("TestFieldProps failed: 'Embedding' not changed to 'float[]'")
	}

	err := go_db_entities.Create(ctx, &TestEntity{})
	if err != nil {
		t.Error(err)
	}
}
