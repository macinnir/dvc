package request

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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

	if valString != "" {
		if val, e = strconv.ParseInt(valString, 10, 64); e == nil {
			return
		}
	}

	val = defaultVal
	return

}

// PrintArgs prints args
func PrintArgs(r *http.Request) {

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
			log.Println("Error: ", e.Error())
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
