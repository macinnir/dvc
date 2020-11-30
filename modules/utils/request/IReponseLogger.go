package request

import (
	"net/http"
)

// IResponseLogger outlines methods for handling logging for API Responses
type IResponseLogger interface {

	// HandleError handles errors returned from the service layer and
	// calls a api error handler to return the corresponding HTTP response
	HandleError(r *Request, w http.ResponseWriter, e error)

	// NotImplemented shows a text response for not implemented method (501)
	NotImplemented(r *Request, w http.ResponseWriter)

	// InternalServerError returns a 500 server error response
	InternalServerError(r *Request, w http.ResponseWriter, e error)

	// NotFound returns a not-found status
	NotFound(r *Request, w http.ResponseWriter)

	// BadRequest returns a bad request status (400)
	BadRequest(r *Request, w http.ResponseWriter, e error)

	// Unauthorized returns an unauthorized status (401)
	Unauthorized(r *Request, w http.ResponseWriter)

	// Forbidden returns a forbidden status (403)
	Forbidden(r *Request, w http.ResponseWriter)

	// NoContent returns a noContent status (204)
	NoContent(r *Request, w http.ResponseWriter)

	// Created returns a created status (201)
	Created(r *Request, w http.ResponseWriter)

	// JSON Returns an ok status with json-encoded body
	JSON(r *Request, w http.ResponseWriter, body interface{})

	// OK Returns an ok status
	OK(r *Request, w http.ResponseWriter)
}
