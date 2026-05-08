package query_test

import (
	"testing"

	"github.com/macinnir/dvc/core/lib/utils/query"
	"github.com/stretchr/testify/assert"
)

func TestEscapeString(t *testing.T) {

	result := query.EscapeString("I'm a string")
	assert.Equal(t, `I\'m a string`, result)

	result = query.EscapeString(`I"m a string`)
	assert.Equal(t, `I\"m a string`, result)

}
