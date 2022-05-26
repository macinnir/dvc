package gen

import (
	"testing"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var structBytes []byte = []byte(`package foo

import (
	foo "bar"
	"fmt"
	"log"
	_ "wtf"
)

// Foo is a foo
type Foo struct {
	A string ` + "`s1:\"s1ValA\" s2:\"s2ValA,optsA\"` // Here are some comments!" + `
	B int64  ` + "`s1:\"s1ValB\" s2:\"s2ValB,optsB\"`" + `
	C bool   ` + "`s1:\"s1ValC\" s2:\"s2ValC,optsC\"` // C Comments" + `
}
`)

var tableB *schema.Table = &schema.Table{
	Name: "School",
	Columns: map[string]*schema.Column{
		"SchoolID":    {Name: "SchoolID", DataType: "int", IsNullable: false},
		"Address":     {Name: "Address", DataType: "varchar", IsNullable: false},
		"Address2":    {Name: "Address2", DataType: "varchar", IsNullable: false},
		"City":        {Name: "City", DataType: "varchar", IsNullable: false},
		"State":       {Name: "State", DataType: "varchar", IsNullable: false},
		"County":      {Name: "County", DataType: "varchar", IsNullable: false},
		"DateCreated": {Name: "County", DataType: "datetime", IsNullable: true},
		"LastUpdated": {Name: "LastUpdated", DataType: "bigint", IsNullable: false},
		"IsDeleted":   {Name: "IsDeleted", DataType: "tinyint", IsNullable: false},
		"Name":        {Name: "Name", DataType: "varchar", IsNullable: false},
		"Zip":         {Name: "Zip", DataType: "varchar", IsNullable: false},
	},
}

var structB []byte = []byte(`package models

import (
	` + NullPackage + `
)

// School represents a School domain object
type School struct {
	SchoolID    int64       ` + "`db:\"SchoolID\" json:\"SchoolID\"`" + `
	Name        string      ` + "`db:\"Name\" json:\"Name\"`" + `
	IsDeleted   int         ` + "`db:\"IsDeleted\" json:\"IsDeleted\"`" + `
	Address     string      ` + "`db:\"Address\" json:\"Address\"`" + `
	Address2    string      ` + "`db:\"Address2\" json:\"Address2\"`" + `
	City        string      ` + "`db:\"City\" json:\"City\"`" + `
	State       string      ` + "`db:\"State\" json:\"State\"`" + `
	Zip         string      ` + "`db:\"Zip\" json:\"Zip\"`" + `
	County      string      ` + "`db:\"County\" json:\"County\"`" + `
	DateCreated null.String ` + "`db:\"DateCreated\" json:\"DateCreated\"`" + `
}
`)

var structBUpdated []byte = []byte(`package models

import (
	` + NullPackage + `
)

// School represents a School domain object
type School struct {
	SchoolID    int64       ` + "`db:\"SchoolID\" json:\"SchoolID\"`" + `
	Name        string      ` + "`db:\"Name\" json:\"Name\"`" + `
	IsDeleted   int         ` + "`db:\"IsDeleted\" json:\"IsDeleted\"`" + `
	Address     string      ` + "`db:\"Address\" json:\"Address\"`" + `
	Address2    string      ` + "`db:\"Address2\" json:\"Address2\"`" + `
	City        string      ` + "`db:\"City\" json:\"City\"`" + `
	State       string      ` + "`db:\"State\" json:\"State\"`" + `
	Zip         string      ` + "`db:\"Zip\" json:\"Zip\"`" + `
	County      string      ` + "`db:\"County\" json:\"County\"`" + `
	DateCreated null.String ` + "`db:\"DateCreated\" json:\"DateCreated\"`" + `
	LastUpdated int64       ` + "`db:\"LastUpdated\" json:\"LastUpdated\"`" + `
}
`)

func TestParseStringToGoStruct(t *testing.T) {

	modelNode, e := parseStringToGoStruct(structBytes)

	require.Nil(t, e)
	assert.Equal(t, "Foo", modelNode.Name)
	assert.Equal(t, "foo", modelNode.Package)
	assert.Equal(t, "Foo is a foo\n", modelNode.Comments)

	// Imports
	assert.Equal(t, 4, modelNode.Imports.Len())
	assert.Equal(t, "foo \"bar\"", modelNode.Imports.Get(0))
	assert.Equal(t, "\"fmt\"", modelNode.Imports.Get(1))
	assert.Equal(t, "\"log\"", modelNode.Imports.Get(2))
	assert.Equal(t, "_ \"wtf\"", modelNode.Imports.Get(3))

	assert.Equal(t, "A", modelNode.Fields.Get(0).Name)
	assert.Equal(t, "string", modelNode.Fields.Get(0).DataType)
	assert.Equal(t, "Here are some comments!\n", modelNode.Fields.Get(0).Comments)
	assert.Equal(t, "s1", modelNode.Fields.Get(0).Tags[0].Name)
	assert.Equal(t, "s1ValA", modelNode.Fields.Get(0).Tags[0].Value)
	assert.Equal(t, "s2", modelNode.Fields.Get(0).Tags[1].Name)
	assert.Equal(t, "s2ValA", modelNode.Fields.Get(0).Tags[1].Value)
	assert.Equal(t, "optsA", modelNode.Fields.Get(0).Tags[1].Options[0])

	assert.Equal(t, "B", modelNode.Fields.Get(1).Name)
	assert.Equal(t, "int64", modelNode.Fields.Get(1).DataType)
	assert.Equal(t, "", modelNode.Fields.Get(1).Comments)
	assert.Equal(t, "s1", modelNode.Fields.Get(1).Tags[0].Name)
	assert.Equal(t, "s1ValB", modelNode.Fields.Get(1).Tags[0].Value)
	assert.Equal(t, "s2", modelNode.Fields.Get(1).Tags[1].Name)
	assert.Equal(t, "s2ValB", modelNode.Fields.Get(1).Tags[1].Value)
	assert.Equal(t, "optsB", modelNode.Fields.Get(1).Tags[1].Options[0])

	assert.Equal(t, "C", modelNode.Fields.Get(2).Name)
}

func TestResolveTableToModel_WithNulls(t *testing.T) {
	table := &schema.Table{
		Columns: map[string]*schema.Column{
			"One":         {Name: "One", DataType: "varchar"},
			"NullableCol": {Name: "NullableCol", DataType: "varchar", IsNullable: true},
		},
	}

	// Model has no null import
	model := lib.NewGoStruct()
	model.Fields = &lib.GoStructFields{}
	model.Fields.Append(&lib.GoStructField{Name: "One", DataType: "string"})

	resolveTableToModel(model, table)

	require.Equal(t, 1, model.Imports.Len())
	assert.Equal(t, NullPackage, model.Imports.Get(0))
	assert.Equal(t, 2, model.Fields.Len())
	assert.Equal(t, "NullableCol", model.Fields.Get(1).Name)
	assert.Equal(t, "null.String", model.Fields.Get(1).DataType)
}

func TestResolveTableToModel_UpdatedType(t *testing.T) {
	table := &schema.Table{
		Columns: map[string]*schema.Column{
			"One":         {Name: "One", DataType: "varchar"},
			"NullableCol": {Name: "NullableCol", DataType: "int", IsNullable: true},
		},
	}

	// Model has no null import
	model := lib.NewGoStruct()
	model.Fields = &lib.GoStructFields{}
	model.Fields.Append(&lib.GoStructField{Name: "One", DataType: "string"})
	model.Fields.Append(&lib.GoStructField{Name: "NullableCol", DataType: "string"})

	resolveTableToModel(model, table)

	require.Equal(t, 1, model.Imports.Len())
	assert.Equal(t, NullPackage, model.Imports.Get(0))
	assert.Equal(t, 2, model.Fields.Len())
	assert.Equal(t, "NullableCol", model.Fields.Get(1).Name)
	assert.Equal(t, "null.Int", model.Fields.Get(1).DataType)
}

func TestResolveTableToModel_RemoveColumn(t *testing.T) {
	table := &schema.Table{
		Columns: map[string]*schema.Column{
			"One": {Name: "One", DataType: "varchar"},
		},
	}

	// Model has no null import
	model := lib.NewGoStruct()
	model.Imports = &lib.GoFileImports{NullPackage}
	model.Fields = &lib.GoStructFields{}
	model.Fields.Append(&lib.GoStructField{Name: "One", DataType: "string"})
	model.Fields.Append(&lib.GoStructField{Name: "Two", DataType: "null.String"})

	resolveTableToModel(model, table)

	require.Equal(t, 0, model.Imports.Len())
	assert.Equal(t, 1, model.Fields.Len())
	assert.Equal(t, "One", model.Fields.Get(0).Name)
}

func TestBuildModelNodeFromTable(t *testing.T) {
	table := &schema.Table{
		Name: "Foo",
		Columns: map[string]*schema.Column{
			"Foo1":    {Name: "Foo1", DataType: "int"},
			"NullFoo": {Name: "NullFoo", DataType: "datetime", IsNullable: true},
		},
	}
	m, e := buildModelNodeFromTable(table)
	require.Nil(t, e)
	assert.Equal(t, "Foo", m.Name)
	require.Equal(t, 1, m.Imports.Len())
	assert.Equal(t, NullPackage, m.Imports.Get(0))
	require.Equal(t, 2, m.Fields.Len())
	assert.Equal(t, "Foo1", m.Fields.Get(0).Name)
	assert.Equal(t, "int64", m.Fields.Get(0).DataType)
	assert.Equal(t, "NullFoo", m.Fields.Get(1).Name)
	assert.Equal(t, "null.String", m.Fields.Get(1).DataType)
	assert.Contains(t, *m.Imports, NullPackage)
}

func TestParseFileNameToModelName(t *testing.T) {

	var tests = []struct {
		file   string
		prefix string
		suffix string
		result string
	}{
		{"IFooDAL.go", "I", "DAL", "Foo"},
		{"Foo.go", "", "", "Foo"},
		{"Foo_test.go", "", "", "Foo"},
		{"FooDAL_test.go", "", "DAL", "Foo"},
	}

	for k := range tests {

		ts := tests[k]

		var result = parseFileNameToModelName(ts.file, ts.prefix, ts.suffix)
		assert.Equal(t, ts.result, result)

	}

}
