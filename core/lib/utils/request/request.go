package request

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// AuthHeaderKey is the name of the Authorization header
const AuthHeaderKey = "Authorization"

// AuthHeaderValuePrefix is the start of the Authorization string (in the authorization header) used to authorize the request
const AuthHeaderValuePrefix = "Bearer "

// Request is a request object
type Request struct {
	Method           string
	Path             string
	Headers          map[string]string
	Params           map[string]string
	Body             string
	ControllerMethod string
	ResponseCode     int64
	Error            string
	UserID           int64
	RootRequest      *http.Request `json:"-"`
	ActionType       int64
}

// NewRequest is a factory method for a request
func NewRequest(r *http.Request) *Request {

	req := &Request{
		Method:      r.Method,
		Path:        r.URL.RequestURI(),
		Headers:     map[string]string{},
		Params:      map[string]string{},
		Body:        "",
		RootRequest: r,
		ActionType:  UserLogActionTypeAPI.Int64(),
	}

	// Headers
	for name := range r.Header {
		if len(r.Header[name]) > 0 {
			req.Headers[name] = r.Header[name][0]
		}
	}

	// URL Params
	req.Params = mux.Vars(r)

	defer r.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(r.Body)

	req.Body = string(bodyBytes)

	return req
}

// ToJSONString returns a json string representation of the request
func (r *Request) ToJSONString() string {

	j, _ := json.MarshalIndent(r, "", "    ")

	return string(j)
}

// ArgInt64 returns an int64 value for an argument in the request named `name`
// If it exists, attempts to part it to a 64-bit integer
// Else, `defaultVal` is returned
func (r *Request) ArgInt64(name string, defaultVal int64) int64 {

	val := int64(0)
	if _, ok := r.Params[name]; !ok {
		return defaultVal
	}

	var e error
	if val, e = strconv.ParseInt(r.Params[name], 10, 64); e != nil {
		return defaultVal
	}

	return val
}

// ArgInt returns an int value for an argument in the request named `name`
// If it exists, attempts to part it to a 64-bit integer
// Else, `defaultVal` is returned
func (r *Request) ArgInt(name string, defaultVal int) int {

	val := int(0)
	if _, ok := r.Params[name]; !ok {
		return defaultVal
	}

	var e error
	if val, e = strconv.Atoi(r.Params[name]); e != nil {
		return defaultVal
	}

	return val
}

// Arg returns a string value for an argument in the request named `name`
// If it exists, returns it as a string
// Else, returns the `defaultVal`
func (r *Request) Arg(name string, defaultVal string) string {

	if _, ok := r.Params[name]; !ok {
		return defaultVal
	}

	return r.Params[name]
}

// Header gets a header by its name
func (r *Request) Header(name string) string {
	if _, ok := r.Headers[name]; ok {
		return r.Headers[name]
	}

	return ""
}

// BodyJSON extracts the json from the body of a post or put request
func (r *Request) BodyJSON(obj interface{}) {
	json.Unmarshal([]byte(r.Body), obj)
}

// AuthKey returns the authorization key from the request header
func (r *Request) AuthKey() string {
	authKey := r.Header(AuthHeaderKey)
	if len(AuthHeaderValuePrefix) >= len(authKey) {
		return ""
	}

	return authKey[len(AuthHeaderValuePrefix):]
}