package request

import "strings"

// HTTPMethodType is the type of a method request
type HTTPMethodType int

const (
	// HTTPMethodUnknown is an unknown method type
	HTTPMethodUnknown HTTPMethodType = iota
	// HTTPMethodGet is the http method 'GET'
	HTTPMethodGet
	// HTTPMethodPost is the http method POST
	HTTPMethodPost
	// HTTPMethodPut is the http method PUT
	HTTPMethodPut
	// HTTPMethodDelete is the http method DELETE
	HTTPMethodDelete
	// HTTPMethodOptions is the http method OPTIONS
	HTTPMethodOptions
	// HTTPMethodHead is the http method HEAD
	HTTPMethodHead
	// HTTPMethodTrace is the http method TRACE
	HTTPMethodTrace
	// HTTPMethodPatch is the http method PATCH
	HTTPMethodPatch
)

// HTTPMethodTypeFromString converts a method from a string to its method type
func HTTPMethodTypeFromString(method string) HTTPMethodType {

	method = strings.ToLower(method)

	switch method {
	case "get":
		return HTTPMethodGet
	case "post":
		return HTTPMethodPost
	case "put":
		return HTTPMethodPut
	case "delete":
		return HTTPMethodDelete
	case "options":
		return HTTPMethodOptions
	case "head":
		return HTTPMethodHead
	case "trace":
		return HTTPMethodTrace
	case "patch":
		return HTTPMethodPatch
	default:
		return HTTPMethodUnknown
	}
}

func HTTPMethodTypeToString(t HTTPMethodType) string {
	switch t {
	case HTTPMethodGet:
		return "GET"
	case HTTPMethodPost:
		return "POST"
	case HTTPMethodPut:
		return "PUT"
	case HTTPMethodDelete:
		return "DELETE"
	case HTTPMethodOptions:
		return "OPTIONS"
	case HTTPMethodHead:
		return "HEAD"
	case HTTPMethodTrace:
		return "TRACE"
	case HTTPMethodPatch:
		return "PATCH"
	default:
		return "UNKNOWN"
	}
}

// Int64 returns the int64 representation of the HTTPMethodType value
func (h HTTPMethodType) Int64() int64 {
	return int64(h)
}
