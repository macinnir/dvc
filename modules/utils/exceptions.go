package utils

import (
	"fmt"
)

// RecordNotFoundError is a Not FOund (404) network response
type RecordNotFoundError struct {
}

func (e RecordNotFoundError) Error() string {
	return "Not Found"
}

// NewRecordNotFoundError returns a RecordNotFoundError object
func NewRecordNotFoundError() RecordNotFoundError {
	return RecordNotFoundError{}
}

// NotAuthorizedError is a NotAuthorized (401) Network Response
type NotAuthorizedError struct {
}

func (e NotAuthorizedError) Error() string {
	return "Not Authorized"
}

// NewNotAuthorizedError returns a new NotAuthorizedError object
func NewNotAuthorizedError() NotAuthorizedError {
	return NotAuthorizedError{}
}

// ForbiddenError - 403
type ForbiddenError struct{}

func (e ForbiddenError) Error() string {
	return fmt.Sprintf("Forbidden")
}

// NewForbiddenError returns a ForbiddenError object
func NewForbiddenError() ForbiddenError {
	return ForbiddenError{}
}

// NewInternalError returns an InternalError object
func NewInternalError(s string) InternalError {
	return InternalError{s}
}

// InternalError represents a network response of Internal Server Error (500)
type InternalError struct {
	s string
}

func (e InternalError) Error() string {
	return e.s
}

// NewArgumentError returns a new ArgumentError object 
func NewArgumentError(s string) ArgumentError {
	return ArgumentError{s}
}

// ArgumentError represents a BadRequest (400) network response
type ArgumentError struct {
	s string
}

func (e ArgumentError) Error() string {
	return e.s
}
