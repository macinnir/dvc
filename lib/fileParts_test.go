package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGoStruct(t *testing.T) {
	newStruct := NewGoStruct()
	assert.NotNil(t, newStruct.Fields)
	assert.NotNil(t, newStruct.Imports)
}

func TestGoStructFields(t *testing.T) {

	fields := &GoStructFields{}
	assert.Equal(t, 0, fields.Len())

	fields.Append(&GoStructField{
		Name: "FirstField",
	})
	assert.Equal(t, 1, fields.Len())

	require.NotNil(t, fields.Get(0))
	assert.Equal(t, "FirstField", fields.Get(0).Name)
}

func TestGoFileImportsToString(t *testing.T) {
	i := &GoFileImports{"_ \"a\"", "\"b\"", "d \"c\""}
	assert.Equal(t, "_ \"a\"", i.Get(0))
	assert.Equal(t, 3, i.Len())
	assert.Equal(t, `import (
	_ "a"
	"b"
	d "c"
)
`, i.ToString())
}

func TestGoStructFieldToString(t *testing.T) {
	// modelNode, e := buildModelNodeFromFile(structBytes)
	field := &GoStructField{
		Name:     "A",
		DataType: "string",
		Tags: []*GoStructFieldTag{
			{Name: "s1", Value: "s1ValA", Options: []string{}},
			{Name: "s2", Value: "s2ValA", Options: []string{"optsA"}},
		},
		Comments: "Here are some comments!",
	}
	// assert.Nil(t, e)
	assert.Equal(t, "A string `s1:\"s1ValA\" s2:\"s2ValA,optsA\"` // Here are some comments!", field.ToString())
}

func TestGoStructFieldTagToString(t *testing.T) {
	// modelNode, e := buildModelNodeFromFile(structBytes)
	tag := &GoStructFieldTag{
		Name:  "foo",
		Value: "bar",
		Options: []string{
			"one",
			"two",
			"three",
		},
	}
	// assert.Nil(t, e)
	assert.Equal(t, "foo:\"bar,one,two,three\"", tag.ToString())
}
