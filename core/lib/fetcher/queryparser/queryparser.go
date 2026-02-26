package queryparser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/macinnir/dvc/core/lib"
)

func ExtractParamFromString(paramString string) (param lib.ControllerRouteParam) {

	// Incase there are parts after the param, split on the closing bracket
	pParts := strings.Split(paramString, "}")
	paramString = pParts[0]

	paramParts := strings.Split(paramString, ":")

	param = lib.ControllerRouteParam{
		Name:    paramParts[0],
		Pattern: paramParts[1],
	}

	param.Type = matchPatternToDataType(param.Pattern)
	return
}

func matchPatternToDataType(pattern string) string {
	if pattern == "[0-9]" || pattern == "[0-9]+" || pattern == "-?[0-9]+" {
		return "int64"
	}

	return "string"
}

func ExtractQueryStringsFromRoutePath(routePath string) (string, []string) {
	if !strings.Contains(routePath, "?") {
		return routePath, []string{}
	}
	subParts := strings.SplitN(routePath, "?", 2)
	return subParts[0], strings.Split(subParts[1], "&")
}

func ExtractQueriesFromRoutePath(routePath string) (string, []lib.ControllerRouteQuery) {

	queries := []lib.ControllerRouteQuery{}

	path, queryStrings := ExtractQueryStringsFromRoutePath(routePath)

	for _, query := range queryStrings {

		if !strings.Contains(query, "=") {
			continue
		}

		o := ParseURLQuerySegment(query)

		queries = append(queries, o)
	}

	return path, queries
}

// ParseURLQuerySegment parses a query segment in a URL query to a lib.ControllerRouteQueryObject
//
//	e.g. foo={foo:[0-9]+}
func ParseURLQuerySegment(pattern string) lib.ControllerRouteQuery {

	queryParts := strings.Split(pattern, "=")

	o := lib.ControllerRouteQuery{
		Name:     queryParts[0],
		ValueRaw: queryParts[1],
	}

	if strings.Contains(o.ValueRaw, ":") {
		queryValueParts := strings.Split(o.ValueRaw, ":")
		// Remove the starting "{"
		o.VariableName = queryValueParts[0][1:]

		// Remove the ending "}"
		o.Pattern = strings.Join(queryValueParts[1:], ":")
		o.Pattern = o.Pattern[0 : len(o.Pattern)-1]

		// Check if the value isn't a constant value
		o.Type = matchPatternToDataType(o.Pattern)

	} else {
		// Try to parse the value as an int64
		// e.g. param=123

		o.VariableName = o.Name
		if _, e := strconv.ParseInt(o.ValueRaw, 10, 64); e != nil {
			o.Type = "string"
		} else {
			o.Type = "int64"
		}
	}

	return o
}

func ParseRouteString(route *lib.ControllerRoute, routeString string) error {
	
	lineParts := strings.SplitN(routeString, " ", 4)
	if len(lineParts) < 4 {
		return fmt.Errorf("invalid route comment `%s`", routeString)
	}
	route.Method = lineParts[2]
	route.Raw = lineParts[3]

	route.Path, route.Queries = ExtractQueriesFromRoutePath(route.Raw)
	for _, q := range route.Queries {
		if q.Type != "int64" && strings.HasPrefix(q.Pattern, "-?") {
			continue
		}
		route.RequiredQueries = append(route.RequiredQueries, q)
	}

	params, _ := ExtractParamsFromRoutePath(route.Path)
	route.Params = append(route.Params, params...)
	return nil
}

func ExtractParamsFromRoutePath(routePath string) (params []lib.ControllerRouteParam, e error) {

	params = []lib.ControllerRouteParam{}

	// Params
	if strings.Contains(routePath, "{") {

		routeParts := strings.Split(routePath, "{")

		for _, p := range routeParts[1:] {

			if !strings.Contains(p, "}") || !strings.Contains(p, ":") {
				continue
			}

			param := ExtractParamFromString(p)

			params = append(params, param)
		}
	}

	return
}
