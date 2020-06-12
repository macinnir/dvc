package gen

import (
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestGenInterface(t *testing.T) {
	methods := []string{"// MyMethod does cool stuff", "MyMethod(string) example.Example"}
	imports := []string{`"github.com/example/example"`}
	b, err := GenInterface("DO NOT EDIT: Auto generated", "pkg", "MyInterface", "MyInterface does cool stuff", methods, imports)
	assert.Nil(t, err, "GenInterface returned an error")

	expected := `// DO NOT EDIT: Auto generated

package pkg

import (
	"github.com/example/example"
)

// MyInterface does cool stuff
type MyInterface interface {
	// MyMethod does cool stuff
	MyMethod(string) example.Example
}
`

	assert.Equal(t, expected, string(b))
}

func TestMakeWithoutInterfaceComment(t *testing.T) {
	methods := []string{"// MyMethod does cool stuff", "MyMethod(string) example.Example"}
	imports := []string{`"github.com/example/example"`}
	b, err := GenInterface("DO NOT EDIT: Auto generated", "pkg", "MyInterface", "", methods, imports)
	assert.Nil(t, err, "GenInterface returned an error")

	expected := `// DO NOT EDIT: Auto generated

package pkg

import (
	"github.com/example/example"
)

type MyInterface interface {
	// MyMethod does cool stuff
	MyMethod(string) example.Example
}
`

	assert.Equal(t, expected, string(b))
}

func TestGenInterfaceMultiLineIfaceComment(t *testing.T) {
	b, err := GenInterface("DO NOT EDIT: Auto generated", "pkg", "MyInterface", "MyInterface does cool stuff.\nWith multi-line comments.", nil, nil)
	assert.Nil(t, err, "GenInterface returned an error:", err)

	expected := `// DO NOT EDIT: Auto generated

package pkg

// MyInterface does cool stuff.
// With multi-line comments.
type MyInterface interface {
}
`

	assert.Equal(t, expected, string(b))
}
