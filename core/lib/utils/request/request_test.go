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

func TestIP_Exists(t *testing.T) {
	request := &Request{
		Headers: map[string]string{
			"X-Forwarded-For": "123.123.123.123",
		},
	}

	assert.Equal(t, "123.123.123.123", request.IP())
}

func TestIP_Empty(t *testing.T) {
	request := &Request{
		Headers: map[string]string{},
	}

	assert.Equal(t, "127.0.0.1", request.IP())
}

func TestIP_Multiple(t *testing.T) {
	request := &Request{
		Headers: map[string]string{
			"X-Forwarded-For": "24.180.127.227, 34.117.111.176",
		},
	}

	assert.Equal(t, "24.180.127.227", request.IP())
}
