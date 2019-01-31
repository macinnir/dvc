package utils

import (
	"fmt"
)

type RecordNotFoundError struct {
	TableName string
	Key       string
}

func (e RecordNotFoundError) Error() string {
	return fmt.Sprintf("No record in table `%s` found at key %s", e.TableName, e.Key)
}

func NewRecordNotFoundError(tableName string, key string) RecordNotFoundError {
	return RecordNotFoundError{
		TableName: tableName,
		Key:       key,
	}
}

// ForbiddenError - 403
type ForbiddenError struct{}

func (e ForbiddenError) Error() string {
	return fmt.Sprintf("Forbidden")
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
