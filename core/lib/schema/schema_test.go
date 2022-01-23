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
		{"[]*Foo", "Foo[]"},
		{"*Foo", "Foo"},
		{"map[int64]*Foo", "{ [key: number]: Foo }"},
		{"*models.Foo", "Foo"},
		{"bytes.Buffer", "any"},
		{"*bytes.Buffer", "any"},
		{"interface{}", "any"},
		{"[]byte", "any"},
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
