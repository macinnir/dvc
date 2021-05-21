package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryIDs(t *testing.T) {
	result := QueryIDs("1,2,3")
	assert.Equal(t, int64(1), result[0])
	assert.Equal(t, int64(2), result[1])
	assert.Equal(t, int64(3), result[2])
}

func TestQueryIDs_NonInts(t *testing.T) {
	result := QueryIDs("a,b,c")
	assert.Len(t, result, 0)
}
