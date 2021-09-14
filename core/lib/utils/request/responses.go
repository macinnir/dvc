package request

import (
	"encoding/json"
	"log"
	"net/http"

	gerrors "github.com/go-errors/errors"
	"github.com/macinnir/dvc/core/lib/utils/errors"
)

// HandleError handles errors returned from the service layer and
// calls a api error handler to return the corresponding HTTP response
func HandleError(r *Request, w http.ResponseWriter, e error) {
	// t := reflect.TypeOf(e)
	switch e.(type) {
	case errors.ArgumentError:
		BadRequest(r, w, e)
	case errors.InternalError:
		InternalServerError(r, w, e)
	case errors.ForbiddenError:
		Forbidden(r, w)
	case errors.RecordNotFoundError:
		NotFound(r, w)
	case errors.NotAuthorizedError:
		Unauthorized(r, w)
	default:
		InternalServerError(r, w, e)
	}
}

// NotImplemented shows a text response for not implemented method (501)
func NotImplemented(r *Request, w http.ResponseWriter) {
	r.ResponseCode = 501
	log.Printf(" HTTP %s %s 501 Not Implemented", r.Method, r.Path)
	w.WriteHeader(http.StatusNotImplemented)
	w.Header().Set("content-type", "text/plain")
}

// InternalServerError returns a 500 server error response
func InternalServerError(r *Request, w http.ResponseWriter, e error) {
	r.ResponseCode = 500
	r.Error = e.Error()

	stackError := gerrors.Wrap(e, 1).ErrorStack()
	// fmt.Println("stackError", stackError)
	log.Printf(" HTTP %s %s 500 INTERNAL SERVER ERROR: %s", r.Method, r.Path, stackError)
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("content-type", "text/plain")
}

// ErrorResponse is the structure of a response that is an error
// @model ErrorResponse
type ErrorResponse struct {
	Status string `json:"status"`
	Detail string `json:"detail"`
}

// NotFound returns a not-found status
func NotFound(r *Request, w http.ResponseWriter) {
	r.ResponseCode = 404
	log.Printf("WAR HTTP %s %s 404 NOT FOUND", r.Method, r.Path)
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("content-type", "text/plain")
}

// BadRequest returns a bad request status (400)
func BadRequest(r *Request, w http.ResponseWriter, e error) {
	log.Printf("WAR HTTP %s %s 400 BAD REQUEST: %s", r.Method, r.Path, e.Error())
	w.WriteHeader(http.StatusBadRequest)
	r.ResponseCode = 400
	r.Error = e.Error()
	errorResponse := ErrorResponse{}
	errorResponse.Status = "400"
	errorResponse.Detail = e.Error()
	JSON(r, w, errorResponse)
}

// Unauthorized returns an unauthorized status (401)
func Unauthorized(r *Request, w http.ResponseWriter) {
	r.ResponseCode = 401
	w.WriteHeader(http.StatusUnauthorized)
	log.Printf("WAR HTTP %s %s 401 NOT AUTHORIZED", r.Method, r.Path)
	w.Header().Set("content-type", "text/plain")
}

// Forbidden returns a forbidden status (403)
func Forbidden(r *Request, w http.ResponseWriter) {
	r.ResponseCode = 403
	log.Printf("WAR HTTP %s %s 403 FORBIDDEN", r.Method, r.Path)
	w.WriteHeader(http.StatusForbidden)
	w.Header().Set("content-type", "text/plain")
}

// NoContent returns a noContent status (204)
func NoContent(r *Request, w http.ResponseWriter) {
	r.ResponseCode = 204
	log.Printf("INF HTTP %s %s 204 NO CONTENT", r.Method, r.Path)
	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("content-type", "text/plain")
}

// Created returns a created status (201)
func Created(r *Request, w http.ResponseWriter) {
	r.ResponseCode = 201
	log.Printf("INF HTTP %s %s 201 CREATED", r.Method, r.Path)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "text/plain")
}

// JSON Returns an ok status with json-encoded body
func JSON(r *Request, w http.ResponseWriter, body interface{}) {
	r.ResponseCode = 200
	log.Printf("INF HTTP %s %s 200 OK", r.Method, r.Path)
	payload, _ := json.Marshal(body)
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func HTML(r *Request, w http.ResponseWriter, content string) {
	r.ResponseCode = 200
	log.Printf("INF HTTP %s %s 200 OK", r.Method, r.Path)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(content))
}

// OK Returns an ok status
func OK(r *Request, w http.ResponseWriter) {
	r.ResponseCode = 200
	log.Printf("INF HTTP %s %s 200 OK", r.Method, r.Path)
	w.Header().Set("content-type", "text/plain")
}
