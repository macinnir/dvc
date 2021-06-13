package gen

import (
	"testing"

	"github.com/macinnir/dvc/core/lib/schema"
	"github.com/stretchr/testify/assert"
)

func TestGenerateTypescriptTypes(t *testing.T) {

	database := &schema.Schema{
		Tables: map[string]*schema.Table{
			"Foo": {
				Name: "Foo",
				Columns: map[string]*schema.Column{
					"FooID": {Name: "FooID", DataType: "int"},
					"Bar":   {Name: "Bar", DataType: "varchar"},
					"Baz":   {Name: "Baz", DataType: "datetime"},
				},
			},
		},
	}

	g, e := GenerateTypescriptTypes(database)

	assert.Nil(t, e)
	assert.Equal(t, "// #genStart \n\ndeclare namespace Models {\n\n\t/**\n\t * Foo\n\t */\n\texport interface Foo{\n\t\tFooID: number;\n\t\tBar: string;\n\t\tBaz: string;\n\t}\n\n}\n// #genEnd\n", g)
}

func TestGenerateTypescriptType(t *testing.T) {

	table := &schema.Table{
		Name: "Foo",
		Columns: map[string]*schema.Column{
			"FooID": {Name: "FooID", DataType: "int"},
			"Bar":   {Name: "Bar", DataType: "varchar"},
			"Baz":   {Name: "Baz", DataType: "datetime"},
		},
	}

	g, e := GenerateTypescriptType(table)

	assert.Nil(t, e)
	assert.Equal(t, "\t/**\n\t * Foo\n\t */\n\texport interface Foo{\n\t\tFooID: number;\n\t\tBar: string;\n\t\tBaz: string;\n\t}\n\n", g)
}
