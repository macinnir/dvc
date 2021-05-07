package dump

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBFielCleanString(t *testing.T) {
	val := `Flier's Quality Water Systems, Inc.`

	expected := `Flier\'s Quality Water Systems, Inc.`

	assert.Equal(t, expected, dbFieldCleanString(val))
}
