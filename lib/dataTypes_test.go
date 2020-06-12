package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataTypeToGoTypeString(t *testing.T) {
	assert.Equal(t, "int", DataTypeToGoTypeString(&Column{DataType: "tinyint"}))
	assert.Equal(t, "int64", DataTypeToGoTypeString(&Column{DataType: "int"}))
	assert.Equal(t, "string", DataTypeToGoTypeString(&Column{DataType: "varchar"}))
	assert.Equal(t, "string", DataTypeToGoTypeString(&Column{DataType: "enum"}))
	assert.Equal(t, "string", DataTypeToGoTypeString(&Column{DataType: "text"}))
	assert.Equal(t, "string", DataTypeToGoTypeString(&Column{DataType: "date"}))
	assert.Equal(t, "string", DataTypeToGoTypeString(&Column{DataType: "datetime"}))
	assert.Equal(t, "string", DataTypeToGoTypeString(&Column{DataType: "char"}))
	assert.Equal(t, "float64", DataTypeToGoTypeString(&Column{DataType: "decimal"}))
	assert.Equal(t, "null.String", DataTypeToGoTypeString(&Column{DataType: "varchar", IsNullable: true}))
	assert.Equal(t, "null.Int", DataTypeToGoTypeString(&Column{DataType: "int", IsNullable: true}))
	assert.Equal(t, "null.Float", DataTypeToGoTypeString(&Column{DataType: "decimal", IsNullable: true}))
}
