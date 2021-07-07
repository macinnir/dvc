package fetcher

import (
	"testing"

	"github.com/macinnir/dvc/core/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractControllerNameFromFileName(t *testing.T) {
	var tests = []struct {
		path string
		name string
	}{
		{"AController.go", "A"},
		{"AController_test.go", ""},
		{"ALongControllerName", ""},
		{"foo/bar/BazController.go", "Baz"},
	}

	for k := range tests {
		assert.Equal(t, tests[k].name, extractControllerNameFromFileName(tests[k].path))
	}
}

func TestParseRouteString(t *testing.T) {
	route := &lib.ControllerRoute{}
	e := parseRouteString(route, "// @route GET /test?one={one:[0-9]+}")

	assert.Nil(t, e)
	assert.Equal(t, "/test", route.Path)
	assert.Equal(t, "GET", route.Method)
	require.Equal(t, 1, len(route.Queries))
	assert.Equal(t, "one", route.Queries[0].Name)
	assert.Equal(t, "[0-9]+", route.Queries[0].Pattern)
	assert.Equal(t, "int64", route.Queries[0].Type)

}

func TestExtractParamsFromRoutePath(t *testing.T) {

	route := "/get/{getID:[0-9]+}/something/{somethingID:[a-zA-Z-]+}/cool/{coolID:[0-9]+}"

	params, e := extractParamsFromRoutePath(route)

	require.Nil(t, e)
	require.Len(t, params, 3)

	assert.Equal(t, "getID", params[0].Name)
	assert.Equal(t, "[0-9]+", params[0].Pattern)
	assert.Equal(t, "int64", params[0].Type)

	assert.Equal(t, "somethingID", params[1].Name)
	assert.Equal(t, "[a-zA-Z-]+", params[1].Pattern)
	assert.Equal(t, "string", params[1].Type)

	assert.Equal(t, "coolID", params[2].Name)
	assert.Equal(t, "[0-9]+", params[2].Pattern)
	assert.Equal(t, "int64", params[2].Type)

}

func TestExtractParamFromString(t *testing.T) {

	tests := []struct {
		p            string
		paramName    string
		paramType    string
		paramPattern string
	}{
		{"fooID:[0-9]+}", "fooID", "int64", "[0-9]+"},
		{"barID:[a-zA-Z]+}", "barID", "string", "[a-zA-Z]+"},
	}

	for _, test := range tests {
		t.Run(test.p, func(t *testing.T) {
			param := extractParamFromString(test.p)
			assert.Equal(t, test.paramName, param.Name)
			assert.Equal(t, test.paramType, param.Type)
			assert.Equal(t, test.paramPattern, param.Pattern)
		})
	}
}
