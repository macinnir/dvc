package apitest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const url = "http://someurl.com"

var at *APITest

func TestMain(t *testing.T) {
	initTest(t)
}

func initTest(t *testing.T) {
	at = NewAPITest(t, url)
}

func TestRandString(t *testing.T) {
	r := RandString(10)
	assert.Equal(t, len(r), 10)
	r2 := RandString(10)
	assert.NotEqual(t, r, r2)
}

func TestSetGetAuthKey(t *testing.T) {
	if at == nil {
		initTest(t)
	}
	at.SetAuthKey("12345")
	assert.Equal(t, "12345", at.GetAuthKey())
}

func TestSetGetStringVal(t *testing.T) {
	if at == nil {
		initTest(t)
	}

	at.SetStringVal("foo", "bar")
	assert.Equal(t, "bar", at.GetStringVal("foo"))
}

func TestSetGetStringVal_ShouldReturnEmptyStringIfNotSet(t *testing.T) {
	if at == nil {
		initTest(t)
	}

	assert.Equal(t, "", at.GetStringVal("foo123"))
}

func TestSetGetIntVal(t *testing.T) {
	if at == nil {
		initTest(t)
	}

	at.SetIntVal("foo", 123)
	assert.Equal(t, int64(123), at.GetIntVal("foo"))
}

func TestSetGetIntVal_ShouldReturnNegativeOneIfNotSet(t *testing.T) {
	if at == nil {
		initTest(t)
	}

	assert.Equal(t, int64(-1), at.GetIntVal("missingNumber"))
}
