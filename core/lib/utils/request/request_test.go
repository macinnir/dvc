package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestArgInt64(t *testing.T) {

	headerNameFoo := "foo"
	headerValueBar := "bar"

	req := &Request{
		Method: "POST",
		Path:   "/foo/1/bar",
		Headers: map[string]string{
			headerNameFoo: headerValueBar,
			AuthHeaderKey: AuthHeaderValuePrefix + "some auth key string",
		},
		Params: map[string]string{
			"foo": "1",
			"baz": "qux",
		},
		Body:       "",
		ActionType: 0,
	}

	assert.Equal(t, int64(1), req.ArgInt64("foo", 0))
	assert.Equal(t, 1, req.ArgInt("foo", 0))
	assert.Equal(t, "1", req.Arg("foo", ""))
	assert.Equal(t, int64(0), req.ArgInt64("bar", 0))
	assert.Equal(t, 0, req.ArgInt("bar", 0))
	assert.Equal(t, "", req.Arg("bar", ""))
	assert.Equal(t, int64(0), req.ArgInt64("baz", 0))
	assert.Equal(t, 0, req.ArgInt("baz", 0))
	assert.Equal(t, "qux", req.Arg("baz", ""))

	assert.Equal(t, headerValueBar, req.Header(headerNameFoo), "Header '%s' exists and should return its value of '%s'", headerNameFoo, headerValueBar)
	assert.Equal(t, "", req.Header("bar"), "Header 'bar' does not exist and should be empty")

	assert.Equal(t, "some auth key string", req.AuthKey(), "Auth key should exist as a header")

	// assert.Equal(t, "...", req.BodyJSON())
}

func TestAuthHeaderEmpty(t *testing.T) {
	request := &Request{}
	assert.Equal(t, "", request.AuthKey())
}
