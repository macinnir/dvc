package utils

import (
	"fmt"
)

type RecordNotFoundError struct {
}

func (e RecordNotFoundError) Error() string {
	return "Not Found"
}

func NewRecordNotFoundError() RecordNotFoundError {
	return RecordNotFoundError{}
}

type NotAuthorizedError struct {
}

func (e NotAuthorizedError) Error() string {
	return "Not Authorized"
}

func NewNotAuthorizedError() NotAuthorizedError {
	return NotAuthorizedError{}
}

// ForbiddenError - 403
type ForbiddenError struct{}

func (e ForbiddenError) Error() string {
	return fmt.Sprintf("Forbidden")
}

func NewForbiddenError() ForbiddenError {
	return ForbiddenError{}
}

// InternalServerError - 500
func NewInternalError(s string) InternalError {
	return InternalError{s}
}

type InternalError struct {
	s string
}

func (e InternalError) Error() string {
	return e.s
}

// BadRequest - 400
func NewArgumentError(s string) ArgumentError {
	return ArgumentError{s}
}

type ArgumentError struct {
	s string
}

func (e ArgumentError) Error() string {
	return e.s
}
