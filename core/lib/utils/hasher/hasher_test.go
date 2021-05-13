package hasher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {

	hasher := NewHasher("This is a salt", 30)

	vals := []int{12345}

	hash, e := hasher.Hash(vals)

	assert.Nil(t, e)
	assert.Equal(t, "V84EOgzxXmqaNgRp52RpJGAvQBYwd7", hash)

	var decodedVals []int

	decodedVals, e = hasher.DecodeHash(hash)

	assert.Nil(t, e)

	for k := range vals {
		assert.Equal(t, vals[k], decodedVals[k])
	}
}

func BenchmarkHash(b *testing.B) {

	hasher := NewHasher("This is a salt", 30)

	val := 100000001

	for n := 0; n < b.N; n++ {

		val += n
		hash, _ := hasher.Hash([]int{val})

		hasher.DecodeHash(hash)

	}

}
