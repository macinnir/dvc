package schema_test

import (
	"fmt"
	"testing"

	"github.com/macinnir/dvc/core/lib/schema"
	"github.com/stretchr/testify/assert"
)

func TestDataTypeToGoTypeString(t *testing.T) {

	tests := []struct {
		expects    string
		dataType   string
		isNullable bool
	}{
		{"int", "tinyint", false},
		{"int64", "bigint", false},
		{"int64", "int", false},
		{"string", "varchar", false},
		{"string", "enum", false},
		{"string", "text", false},
		{"string", "date", false},
		{"string", "datetime", false},
		{"string", "char", false},
		{"float64", "decimal", false},
		{"null.String", "varchar", true},
		{"null.Int", "int", true},
		{"null.Float", "decimal", true},
	}

	for k := range tests {
		assert.Equal(t, tests[k].expects, schema.DataTypeToGoTypeString(&schema.Column{
			DataType:   tests[k].dataType,
			IsNullable: tests[k].isNullable,
		}), fmt.Sprintf("%d -> %s", k, tests[k].dataType))
	}
}

func TestGoTypeToTypescriptString(t *testing.T) {

	tests := []struct {
		GoType         string
		TypescriptType string
	}{
		{"int64", "number"},
		{"int", "number"},
		{"float64", "number"},
		{"string", "string"},
		{"null.String", "string"},
		{"bool", "boolean"},
		{"[]string", "string[]"},
		{"[]int", "number[]"},
		{"[]int64", "number[]"},
		{"map[int]int", "{ [key: number]: number }"},
		{"map[int64]int", "{ [key: number]: number }"},
		{"map[int]int64", "{ [key: number]: number }"},
		{"map[int64]int64", "{ [key: number]: number }"},
		{"map[int]int", "{ [key: number]: number }"},
		{"map[float64]int", "{ [key: number]: number }"},
		{"map[float64]int64", "{ [key: number]: number }"},
		{"map[float64]string", "{ [key: number]: string }"},
		{"map[string]string", "{ [key: string]: string }"},
		{"map[string][]int", "{ [key: string]: number[] }"},
		{"[]*Foo", "Foo[]"},
		{"[]*models.Foo", "Foo[]"},
		{"*Foo", "Foo"},
		{"map[int64]*Foo", "{ [key: number]: Foo }"},
		{"*models.Foo", "Foo"},

		{"*dtos.Foo", "Foo"},
		{"dtos.Foo", "Foo"},
		{"[]dtos.Foo", "Foo[]"},
		{"[]*dtos.Foo", "Foo[]"},

		{"*appdtos.Foo", "Foo"},
		{"appdtos.Foo", "Foo"},
		{"[]appdtos.Foo", "Foo[]"},
		{"[]*appdtos.Foo", "Foo[]"},

		{"*aggregates.Foo", "Foo"},
		{"aggregates.Foo", "Foo"},
		{"[]aggregates.Foo", "Foo[]"},
		{"[]*aggregates.Foo", "Foo[]"},

		{"bytes.Buffer", "any"},
		{"*bytes.Buffer", "any"},
		{"interface{}", "any"},
		{"[]byte", "any"},

		{"[][]string", "string[][]"},
	}

	for k := range tests {
		assert.Equal(t, tests[k].TypescriptType, schema.GoTypeToTypescriptString(tests[k].GoType), tests[k].GoType)
	}

}

func TestGoTypeToTypescriptDefault(t *testing.T) {

	tests := []struct {
		GoType         string
		TypescriptType string
	}{
		{"int64", "0"},
		{"int", "0"},
		{"float64", "0"},
		{"string", "''"},
		{"null.String", "''"},
		{"bool", "false"},
		{"[]string", "[]"},
		{"[]int", "[]"},
		{"[]int64", "[]"},
		{"map[int]int", "[]"},
		{"map[int64]int", "[]"},
		{"map[int]int64", "[]"},
		{"map[int64]int64", "[]"},
		{"map[int]int", "[]"},
		{"map[float64]int", "[]"},
		{"map[float64]int64", "[]"},
		{"map[float64]string", "[]"},
		{"map[string]string", "{}"},
		{"*Foo", "newFoo()"},
		{"[]*Foo", "[]"},
		{"map[int64]*Foo", "[]"},
		{"bytes.Buffer", "null"},
		{"*bytes.Buffer", "null"},
		{"interface{}", "null"},
		{"[]byte", "null"},

		{"*dtos.Foo", "newFoo()"},
		{"dtos.Foo", "newFoo()"},
		{"[]dtos.Foo", "[]"},
		{"[]*dtos.Foo", "[]"},

		{"*appdtos.Foo", "newFoo()"},
		{"appdtos.Foo", "newFoo()"},
		{"[]appdtos.Foo", "[]"},
		{"[]*appdtos.Foo", "[]"},
	}

	for k := range tests {
		assert.Equal(t, tests[k].TypescriptType, schema.GoTypeToTypescriptDefault(tests[k].GoType), tests[k].GoType)
	}

}

func TestExtractBaseGoType(t *testing.T) {

	tests := []struct {
		GoType   string
		BaseType string
	}{
		{"", ""},
		{"*", ""},
		{"int64", "int64"},
		{"int", "int"},
		{"float64", "float64"},
		{"string", "string"},
		{"null.String", "null.String"},
		{"bool", "bool"},
		{"[]string", "string"},
		{"[]int", "int"},
		{"[]int64", "int64"},
		{"map[int]int", "int"},
		{"map[int64]int", "int"},
		{"map[int]int64", "int64"},
		{"map[int64]int64", "int64"},
		{"map[int]int", "int"},
		{"map[float64]int", "int"},
		{"map[float64]int64", "int64"},
		{"map[float64]string", "string"},
		{"map[string]string", "string"},
		{"*Foo", "Foo"},
		{"[]*Foo", "Foo"},
		{"map[int64]*Foo", "Foo"},
	}

	for k := range tests {
		assert.Equal(t, tests[k].BaseType, schema.ExtractBaseGoType(tests[k].GoType), tests[k].GoType)
	}

}

func TestIsGoTypeBaseType(t *testing.T) {

	tests := []struct {
		GoType string
		IsBase bool
	}{
		{"int64", true},
		{"int", true},
		{"float64", true},
		{"string", true},
		{"null.String", true},
		{"null.Float", true},
		{"bool", true},
		{"[]string", false},
		{"[]int", false},
		{"[]int64", false},
		{"map[int]int", false},
		{"map[int64]int", false},
		{"map[int]int64", false},
		{"map[int64]int64", false},
		{"map[int]int", false},
		{"map[float64]int", false},
		{"map[float64]int64", false},
		{"map[float64]string", false},
		{"map[string]string", false},
		{"*Foo", false},
		{"[]*Foo", false},
		{"map[int64]*Foo", false},
	}

	for k := range tests {
		assert.Equal(t, tests[k].IsBase, schema.IsGoTypeBaseType(tests[k].GoType), tests[k].GoType)
	}

}
