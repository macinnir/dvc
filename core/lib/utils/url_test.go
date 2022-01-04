package utils_test

import (
	"testing"

	"github.com/macinnir/dvc/core/lib/utils"
	"github.com/stretchr/testify/assert"
)

func TestQueryIDs(t *testing.T) {
	result := utils.QueryIDs("1,2,3")
	assert.Equal(t, int64(1), result[0])
	assert.Equal(t, int64(2), result[1])
	assert.Equal(t, int64(3), result[2])
}

func TestQueryIDs_NonInts(t *testing.T) {
	result := utils.QueryIDs("a,b,c")
	assert.Len(t, result, 0)
}

func TestQueryStrings_NoComma(t *testing.T) {
	result := utils.QueryStrings("a")
	assert.Len(t, result, 1)
	assert.Equal(t, result[0], "a")
}
