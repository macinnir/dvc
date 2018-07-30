package gen

import (
	"github.com/macinnir/dvc/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateGoModel(t *testing.T) {

	var e error
	model := ""

	g := &Gen{}

	table := &types.Table{
		Name:    "Foo",
		Columns: map[string]*types.Column{},
	}

	table.Columns["Foo"] = &types.Column{
		Type:     "int(10)",
		DataType: "int",
		Name:     "FooID",
		Position: 1,
	}

	table.Columns["Bar"] = &types.Column{
		Type:       "varchar(32)",
		DataType:   "varchar",
		Name:       "Bar",
		IsNullable: true,
		Position:   2,
	}

	model, e = g.GenerateGoModel(table, []string{})

	assert.Nil(t, e)
	assert.Equal(t, `// #genStart

package models

import "gopkg.in/guregu/null.v3"

// Foo represents a Foo model
type Foo struct {
	FooID int64 `+"`json:\"FooID\"`"+`
	Bar null.String `+"`json:\"Bar\"`"+`
}

// #genEnd`, model)

}
