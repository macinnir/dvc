package errors

import (
	"errors"
	"fmt"
)

// NotAuthorizedError is an error thrown when an action is not authorized for the
// current user
type NotAuthorizedError struct {
}

// Error matches the error interface for NotAuthorizedError
func (e NotAuthorizedError) Error() string {
	return "Not Authorized"
}

// NewNotAuthorizedError returns a NotAuthorizedError
func NewNotAuthorizedError() NotAuthorizedError {
	return NotAuthorizedError{}
}

// ForbiddenError - 403
type ForbiddenError struct{}

// Error matches the error interface for ForbiddenError
func (e ForbiddenError) Error() string {
	return fmt.Sprintf("Forbidden")
}

// NewForbiddenError returns a ForbiddenError
func NewForbiddenError() ForbiddenError {
	return ForbiddenError{}
}

// NewInternalError returns an InternalError
// InternalServerError - 500
func NewInternalError(s string) InternalError {
	return InternalError{s}
}

// InternalError is an error thrown when an action causes an internal
// system error
type InternalError struct {
	s string
}

// Error matches the error interface for InternalError
func (e InternalError) Error() string {
	return e.s
}

// NewArgumentError returns a new ArgumentError
// BadRequest - 400
func NewArgumentError(s string) ArgumentError {
	return ArgumentError{s}
}

// ArgumentError is an error thrown when a method does not have
// the proper input to perform the intended action
type ArgumentError struct {
	s string
}

// Error matches the error interface for ArgumentError
func (e ArgumentError) Error() string {
	return e.s
}

// RecordNotFoundError is an error that occurs when a record is not found in the database
type RecordNotFoundError struct {
}

func (e RecordNotFoundError) Error() string {
	return "Not Found"
}

// NewRecordNotFoundError returns a RecordNotFoundError
func NewRecordNotFoundError() RecordNotFoundError {
	return RecordNotFoundError{}
}

// NewError returns a generic new error
func NewError(text string) error {
	return errors.New(text)
}
