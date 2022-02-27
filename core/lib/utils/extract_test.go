package utils_test

import (
	"testing"

	"github.com/macinnir/dvc/core/lib/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractIDs(t *testing.T) {

	o := []struct {
		id int64
	}{
		{1},
		{1},
		{1},
		{1},
		{1},
		{1},
		{1},
		{1},
		{2},
		{2},
		{2},
		{2},
		{2},
		{2},
		{2},
		{3},
		{3},
		{3},
		{3},
		{3},
		{3},
		{3},
	}

	result := utils.ExtractIDs(len(o), func(index int) int64 { return o[index].id })
	assert.Equal(t, 3, len(result))
	assert.Equal(t, int64(1), result[0])
	assert.Equal(t, int64(2), result[1])
	assert.Equal(t, int64(3), result[2])
}

func TestExtractIndex(t *testing.T) {

	o := []struct {
		id int64
	}{
		{1}, // 0
		{1},
		{1},
		{1},
		{1},
		{1},
		{1},
		{1},
		{2}, // 8
		{2},
		{2},
		{2},
		{2},
		{2},
		{2},
		{3}, // 15
		{3},
		{3},
		{3},
		{3},
		{3},
		{3},
	}

	result := utils.ExtractIndex(len(o), func(index int) int64 { return o[index].id })
	assert.Equal(t, 3, len(result))
	assert.Equal(t, 0, result[1])
	assert.Equal(t, 8, result[2])
	assert.Equal(t, 15, result[3])
}

func TestBatchIDs(t *testing.T) {

	ids := []int64{
		1,
		2,
		3,
		4,
		5,
		6,
		7,
		8,
		9,
		10,
		11,
		12,
		13,
		14,
		15,
		16,
	}

	batches := utils.BatchIDs(ids, 5)

	require.Equal(t, len(batches), 4)
	assert.Equal(t, len(batches[0]), 5)
	assert.Equal(t, len(batches[1]), 5)
	assert.Equal(t, len(batches[2]), 5)
	assert.Equal(t, len(batches[3]), 1)

}
