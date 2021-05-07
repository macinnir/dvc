package compare

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObjectsAreSame(t *testing.T) {

	a := struct {
		a string
	}{
		"astring",
	}

	b := struct {
		b string
	}{
		"bstring",
	}

	result, e := objectsAreSame(a, b)

	assert.Nil(t, e)
	assert.True(t, result)

}
