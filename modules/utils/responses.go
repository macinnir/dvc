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

func NotImplementedHandler(w http.ResponseWriter, r *http.Request) {
	NotImplemented(w)
}

// InternalServerError returns a 500 server error response
func InternalServerError(w http.ResponseWriter, e error) {
	log.Printf("ERROR 500: %s", e.Error())
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
func NotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("content-type", "text/plain")
	return
}

// BadRequest returns a bad request status (400)
func BadRequest(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusBadRequest)
	errorResponse := ErrorResponse{}
	errorResponse.Status = "400"
	errorResponse.Detail = err
	log.Printf("BAD REQUEST (400): %s", err)
	JSON(w, errorResponse)
	return
}

// Unauthorized returns an unauthorized status (401)
func Unauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	log.Println("NOT AUTHORIZED (401)")
	w.Header().Set("content-type", "text/plain")
}

// Forbidden returns a forbidden status (403)
func Forbidden(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	log.Println("FORBIDDEN (403)")
	w.Header().Set("content-type", "text/plain")
}

// NoContent returns a noContent status (204)
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("content-type", "text/plain")
}

// Created returns a created status (201)
func Created(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "text/plain")
}

// JSON Returns an ok status with json-encoded body
func JSON(w http.ResponseWriter, body interface{}) {
	payload, _ := json.Marshal(body)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
}

// OK Returns an ok status
func OK(w http.ResponseWriter) {
	w.Header().Set("content-type", "text/plain")
}

// GetBodyJSON extracts the json from the body of a post or put request
func GetBodyJSON(r *http.Request, obj interface{}) (e error) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	e = decoder.Decode(obj)
	return
}

// UrlVarString returns a string url parameter or the default if not found
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
