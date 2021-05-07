package nonce

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateNonceFromString(t *testing.T) {
	input := "AnotherString"
	result := GenerateNonceFromString(input)
	assert.NotEqual(t, input, result)
}

func TestShortNonceWithPadding(t *testing.T) {
	result := ShortNonceWithPadding(1, 5)
	assert.Len(t, result, 64)
}

func TestBuildShortNonceWithPadding(t *testing.T) {
	result := BuildShortNonceWithPadding(1, "BaseToken")
	assert.Equal(t, "1000000000000000000000000000000000000000000000000000000BaseToken", result)
}
