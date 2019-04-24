package apitest

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRandString(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	r := RandString(10)
	assert.Equal(t, len(r), 10)
	r2 := RandString(10)
	assert.NotEqual(t, r, r2)
}

func TestRandInt(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	r := RandInt(10)
	assert.Greater(t, int64(11), r)
}
