package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmail(t *testing.T) {
	str1 := "ç$€§/az@gmail.com"
	str2 := "abcd@gmail_yahoo.com"
	str3 := "abcd@gmail-yahoo.com"
	str4 := "abcd@gmailyahoo"
	str5 := "abcd@gmail.yahoo"

	assert.False(t, Email(str1))
	assert.False(t, Email(str2))
	assert.True(t, Email(str3))
	assert.True(t, Email(str4))
	assert.True(t, Email(str5))
}
