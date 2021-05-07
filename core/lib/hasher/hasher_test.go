package hasher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {

	vals := []int{12345, 23456}

	hash, e := Hash(vals)

	assert.Nil(t, e)
	assert.Equal(t, "8Znk53NvQyGmkgxfXPbAVPqwxezm6M", hash)

	var decodedVals []int

	decodedVals, e = DecodeHash(hash)

	assert.Nil(t, e)

	assert.Equal(t, vals[0], decodedVals[0])
	assert.Equal(t, vals[1], decodedVals[1])

}
