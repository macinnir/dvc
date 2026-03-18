package model

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
	` + lib.NullPackage + `
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
	` + lib.NullPackage + `
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
	assert.Equal(t, lib.NullPackage, m.Imports.Get(0))
	require.Equal(t, 2, m.Fields.Len())
	assert.Equal(t, "Foo1", m.Fields.Get(0).Name)
	assert.Equal(t, "int64", m.Fields.Get(0).DataType)
	assert.Equal(t, "NullFoo", m.Fields.Get(1).Name)
	assert.Equal(t, "null.String", m.Fields.Get(1).DataType)
	assert.Contains(t, *m.Imports, lib.NullPackage)
}
