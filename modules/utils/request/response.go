package request

import (
	"joc-rfq-api/core/utils/errors"
	"net/http"
)

// Response is the API response handler
type Response struct {
	// userLogDAL    dal.IUserLogDAL
	// userDeviceDAL dal.IUserDeviceDAL
	// auth          integrations.IAuth
	// userDAL       dal.IUserDAL
}

// NewResponse returns a new response object
func NewResponse() *Response {

	return &Response{}

}

// HandleError handles errors returned from the service layer and
// calls a api error handler to return the corresponding HTTP response
func (s *Response) HandleError(r *Request, w http.ResponseWriter, e error) {
	// t := reflect.TypeOf(e)
	switch e.(type) {
	case errors.ArgumentError:
		s.BadRequest(r, w, e)
	case errors.InternalError:
		s.InternalServerError(r, w, e)
	case errors.ForbiddenError:
		s.Forbidden(r, w)
	case errors.RecordNotFoundError:
		s.NotFound(r, w)
	case errors.NotAuthorizedError:
		s.Unauthorized(r, w)
	default:
		s.InternalServerError(r, w, e)
	}
}

// NotImplemented shows a text response for not implemented method (501)
func (s *Response) NotImplemented(r *Request, w http.ResponseWriter) {
	NotImplemented(r, w)
}

// InternalServerError returns a 500 server error response
func (s *Response) InternalServerError(r *Request, w http.ResponseWriter, e error) {
	InternalServerError(r, w, e)
}

// NotFound returns a not-found status
func (s *Response) NotFound(r *Request, w http.ResponseWriter) {
	NotFound(r, w)
}

// BadRequest returns a bad request status (400)
func (s *Response) BadRequest(r *Request, w http.ResponseWriter, e error) {
	BadRequest(r, w, e)
}

// Unauthorized returns an unauthorized status (401)
func (s *Response) Unauthorized(r *Request, w http.ResponseWriter) {
	Unauthorized(r, w)
}

// Forbidden returns a forbidden status (403)
func (s *Response) Forbidden(r *Request, w http.ResponseWriter) {
	Forbidden(r, w)
}

// NoContent returns a noContent status (204)
func (s *Response) NoContent(r *Request, w http.ResponseWriter) {
	NoContent(r, w)
}

// Created returns a created status (201)
func (s *Response) Created(r *Request, w http.ResponseWriter) {
	Created(r, w)
}

// JSON Returns an ok status with json-encoded body
func (s *Response) JSON(r *Request, w http.ResponseWriter, body interface{}) {
	JSON(r, w, body)
}

// OK Returns an ok status
func (s *Response) OK(r *Request, w http.ResponseWriter) {
	OK(r, w)
}
