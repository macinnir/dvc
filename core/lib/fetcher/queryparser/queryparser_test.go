package queryparser_test

import (
	"testing"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/fetcher/queryparser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRouteString(t *testing.T) {
	route := &lib.ControllerRoute{}
	e := queryparser.ParseRouteString(route, "// @route GET /test?one={one:[0-9]+}")

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

	params, e := queryparser.ExtractParamsFromRoutePath(route)

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
		{"bazID:-?[0-9]+}", "bazID", "int64", "-?[0-9]+"},
	}

	for _, test := range tests {
		t.Run(test.p, func(t *testing.T) {
			param := queryparser.ExtractParamFromString(test.p)
			assert.Equal(t, test.paramName, param.Name)
			assert.Equal(t, test.paramType, param.Type)
			assert.Equal(t, test.paramPattern, param.Pattern)
		})
	}
}

func TestExtractQueryStringsFromRoutePath(t *testing.T) {
	rawPath := "/scrubberURls?active={active:-?[0-9]+}&frequency={frequency:[0-9-]+}&status={status:[0-9-]+}&next={next:[0-9]+}&flagged={flagged:[0-9-]+}&folder={folder:[0-9]+}&assignedTo={assignedTo:[0-9-]+}&page={page:[0-9]+}&limit={limit:[0-9]+}&orderBy={orderBy:[a-zA-Z]+}&orderDir={orderDir:[a-zA-Z]+}"

	path, queryStrings := queryparser.ExtractQueryStringsFromRoutePath(rawPath)
	assert.Equal(t, "/scrubberURls", path)
	assert.Equal(t, 11, len(queryStrings))
	assert.Equal(t, "active={active:-?[0-9]+}", queryStrings[0])

}

func TestExtractQueriesFromRoutePath(t *testing.T) {

	rawPath := "/scrubberURls?active={active:-?[0-9]+}&frequency={frequency:[0-9-]+}&status={status:[0-9-]+}&next={next:[0-9]+}&flagged={flagged:[0-9-]+}&folder={folder:[0-9]+}&assignedTo={assignedTo:[0-9-]+}&page={page:[0-9]+}&limit={limit:[0-9]+}&orderBy={orderBy:[a-zA-Z]+}&orderDir={orderDir:[a-zA-Z]+}"

	path, queries := queryparser.ExtractQueriesFromRoutePath(rawPath)

	assert.Equal(t, path, "/scrubberURls")
	assert.Equal(t, 11, len(queries))
	assert.Equal(t, "active", queries[0].Name)
	assert.Equal(t, "-?[0-9]+", queries[0].Pattern)
	assert.Equal(t, "int64", queries[0].Type)
}

func TestParseQuery(t *testing.T) {

	tests := []struct {
		p            string
		paramName    string
		paramType    string
		paramPattern string
	}{
		{"fooID={fooID:[0-9]+}", "fooID", "int64", "[0-9]+"},
		{"barID={barID:[a-zA-Z]+}", "barID", "string", "[a-zA-Z]+"},
		{"bazID={bazID:-?[0-9]+}", "bazID", "int64", "-?[0-9]+"},
		{"quuxID=123", "quuxID", "int64", ""},
		{"quux=foo", "quux", "string", ""},
	}

	for _, test := range tests {
		t.Run(test.p, func(t *testing.T) {
			param := queryparser.ParseURLQuerySegment(test.p)
			assert.Equal(t, test.paramName, param.Name)
			assert.Equal(t, test.paramType, param.Type)
			assert.Equal(t, test.paramPattern, param.Pattern)
		})
	}

}
