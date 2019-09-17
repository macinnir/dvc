package utils

import (
	"encoding/json"

	"github.com/gorilla/mux"

	"log"
	"net/http"
	"strconv"
)

// NotImplemented shows a text response for not implemented method (501)
func NotImplemented(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Header().Set("content-type", "text/plain")
	return
}

// NotImplementedHandler return a new NotImplemented object
func NotImplementedHandler(w http.ResponseWriter, r *http.Request) {
	NotImplemented(w)
}

// InternalServerError returns a 500 server error response
func InternalServerError(r *http.Request, w http.ResponseWriter, e error) {
	log.Printf("ERR HTTP %s %s 500 INTERNAL SERVER ERROR: %s", r.Method, r.RequestURI, e.Error())
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("content-type", "text/plain")
	return
}

// ErrorResponse is the structure of a response that is an error
// @model ErrorResponse
type ErrorResponse struct {
	Status string `json:"status"`
	Detail string `json:"detail"`
}

// NotFound returns a not-found status
func NotFound(r *http.Request, w http.ResponseWriter) {
	log.Printf("WAR HTTP %s %s 404 NOT FOUND", r.Method, r.RequestURI)
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("content-type", "text/plain")
	return
}

// BadRequest returns a bad request status (400)
func BadRequest(r *http.Request, w http.ResponseWriter, e error) {
	log.Printf("WAR HTTP %s %s 400 BAD REQUEST: %s", r.Method, r.RequestURI, e.Error())

	w.WriteHeader(http.StatusBadRequest)
	errorResponse := ErrorResponse{}
	errorResponse.Status = "400"
	errorResponse.Detail = e.Error()
	JSON(r, w, errorResponse)
	return
}

// Unauthorized returns an unauthorized status (401)
func Unauthorized(r *http.Request, w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	log.Printf("WAR HTTP %s %s 401 NOT AUTHORIZED", r.Method, r.RequestURI)
	w.Header().Set("content-type", "text/plain")
}

// Forbidden returns a forbidden status (403)
func Forbidden(r *http.Request, w http.ResponseWriter) {
	log.Printf("WAR HTTP %s %s 403 FORBIDDEN", r.Method, r.RequestURI)
	w.WriteHeader(http.StatusForbidden)
	w.Header().Set("content-type", "text/plain")
}

// NoContent returns a noContent status (204)
func NoContent(r *http.Request, w http.ResponseWriter) {
	log.Printf("INF HTTP %s %s 204 NO CONTENT", r.Method, r.RequestURI)

	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("content-type", "text/plain")
}

// Created returns a created status (201)
func Created(r *http.Request, w http.ResponseWriter) {
	log.Printf("INF HTTP %s %s 201 CREATED", r.Method, r.RequestURI)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "text/plain")
}

// JSON Returns an ok status with json-encoded body
func JSON(r *http.Request, w http.ResponseWriter, body interface{}) {
	log.Printf("INF HTTP %s %s 200 OK", r.Method, r.RequestURI)
	payload, _ := json.Marshal(body)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
}

// OK Returns an ok status
func OK(r *http.Request, w http.ResponseWriter) {
	log.Printf("INF HTTP %s %s 200 OK", r.Method, r.RequestURI)
	w.Header().Set("content-type", "text/plain")
}

// GetBodyJSON extracts the json from the body of a post or put request
func GetBodyJSON(r *http.Request, obj interface{}) (e error) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	e = decoder.Decode(obj)
	return
}

// URLParamString returns a string url parameter or the default if not found
func URLParamString(r *http.Request, name string, defaultVal string) (val string) {
	var ok bool

	vars := mux.Vars(r)

	if val, ok = vars[name]; !ok {
		val = defaultVal
	}

	return
}

// URLParamInt64 returns an int64 value from a url parameter
func URLParamInt64(r *http.Request, name string, defaultVal int64) (val int64) {

	var e error

	valString := URLParamString(r, name, "")

	log.Printf("URLParamInt64: %s\n", valString)

	if valString != "" {
		if val, e = strconv.ParseInt(valString, 10, 64); e == nil {
			return
		}
	}

	val = defaultVal
	return

}

// QueryArgInt checks the incoming request `r` for a query argument named `name`
// and if it exists, attempts to parse it to an integer
// If the argument does not exist, the value `defaultVal` is returned
func QueryArgInt(r *http.Request, name string, defaultVal int) (val int) {

	var e error
	val = 0
	stringVal := r.URL.Query().Get(name)

	if len(stringVal) > 0 {

		val, e = strconv.Atoi(stringVal)

		if e != nil {
			val = defaultVal
			return
		}

		return
	}

	val = defaultVal

	return
}

// QueryArgInt64 checks the incoming request `r` for a query argument named `name`
// and if it exists, attempts to parse it to an 64-bit integer
// If the argument does not exist, the value `defaultVal` is returned
func QueryArgInt64(r *http.Request, name string, defaultVal int64) (val int64) {

	var e error

	val = 0
	stringVal := r.URL.Query().Get(name)

	if len(stringVal) > 0 {

		val, e = strconv.ParseInt(stringVal, 10, 64)

		if e != nil {
			val = defaultVal
			return
		}

		return
	}

	val = defaultVal

	return
}

// QueryArgString checks the incoming request `r` for a query argument named `name`
// and if it exists, returns it
// Else, it returns `defaultVal`
func QueryArgString(r *http.Request, name string, defaultVal string) (val string) {

	stringVal := r.URL.Query().Get(name)

	if len(stringVal) > 0 {
		val = stringVal
		return
	}

	val = defaultVal
	return
}

// APIErrorHandler handles errors returned from the service layer and
// calls a api error handler to return the corresponding HTTP response
func APIErrorHandler(r *http.Request, w http.ResponseWriter, e error) {
	// t := reflect.TypeOf(e)
	switch e.(type) {
	case ArgumentError:
		BadRequest(r, w, e)
	case InternalError:
		InternalServerError(r, w, e)
	case ForbiddenError:
		Forbidden(r, w)
	case RecordNotFoundError:
		NotFound(r, w)
	case NotAuthorizedError:
		Unauthorized(r, w)
	default:
		InternalServerError(r, w, e)
	}
}
