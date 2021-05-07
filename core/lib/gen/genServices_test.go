package gen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsGeneratableFile(t *testing.T) {

	// lowercase
	result := isGeneratableFile("foo.go")
	assert.False(t, result)

	// Uppercase
	result = isGeneratableFile("Foo.go")
	assert.True(t, result)

	// No Extension
	result = isGeneratableFile("Foo")
	assert.False(t, result)

	// No Extension lowercase
	result = isGeneratableFile("foo")
	assert.False(t, result)

	// Test
	result = isGeneratableFile("Foo_test.go")
	assert.False(t, result)

	// Camel case
	result = isGeneratableFile("mockRepos.go")
	assert.False(t, result)

}
