package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildExecManyChunks(t *testing.T) {

	tests := []struct {
		chunkSize int
		stmts     []string
		results   []int
	}{
		{
			5,
			[]string{
				"stmt 1",
				"stmt 2",
				"stmt 3",
				"stmt 4",
				"stmt 5",
				"stmt 6",
				"stmt 7",
				"stmt 8",
				"stmt 9",
				"stmt 10",
				"stmt 11",
				"stmt 12",
				"stmt 13",
				"stmt 14",
				"stmt 15",
				"stmt 16",
				"stmt 17",
				"stmt 18",
				"stmt 19",
				"stmt 20",
				"stmt 21",
				"stmt 22",
			},
			[]int{
				5, 5, 5, 5, 2,
			},
		},
		{5, []string{"stmt 1"}, []int{1}},
		{5, []string{}, []int{0}},
	}

	for k := range tests {
		result := buildExecManyChunks(tests[k].stmts, tests[k].chunkSize)
		for l := range tests[k].results {
			assert.Len(t, result[l], tests[k].results[l])
		}
	}

}
