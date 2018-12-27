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
