package query

import (
	"fmt"
	"log"
)

// Min initiates a count statement
// Count(fieldName string)
func Min(fieldName string) *MinQuery {
	s := &MinQuery{fieldName: fieldName}
	return s
}

// MinQuery is the containing struct for functionality to build count statements on domain objects
type MinQuery struct {
	fieldName string
	Object    DomainObject
	where     []IQueryPart
}

// Where is sugar for the And(queryParts...) method
func (s *MinQuery) Where(queryParts ...IQueryPart) *MinQuery {
	return s.And(queryParts...)
}

// And adds an "and group" to the query
func (s *MinQuery) And(queryParts ...IQueryPart) *MinQuery {

	l := len(queryParts)

	if len(s.where) > 0 {
		s.where = append(s.where, And{})
	}

	group := false
	if l > 0 {
		group = true
		s.where = append(s.where, OpenParenthesis{})
	}

	for _, queryPart := range queryParts {
		s.where = append(s.where, queryPart)
	}

	if group {
		s.where = append(s.where, CloseParenthesis{})
	}
	return s
}

// Or adds an "or group" to the select query
func (s *MinQuery) Or(queryParts ...IQueryPart) *MinQuery {

	l := len(queryParts)

	if len(s.where) > 0 {
		s.where = append(s.where, Or{})
	}

	group := false
	if l > 0 {
		group = true
		s.where = append(s.where, OpenParenthesis{})
	}

	for _, queryPart := range queryParts {
		s.where = append(s.where, queryPart)
	}

	if group {
		s.where = append(s.where, CloseParenthesis{})
	}
	return s
}

// ToSQL builds a sql select statement
func (s *MinQuery) ToSQL() (sql string, args []interface{}) {

	sql, args = buildWhereClauseString(s.where)

	sql = fmt.Sprintf("SELECT MIN(%s) FROM %s %s", escapeField(s.fieldName), escapeField(s.Object.GetName()), sql)

	log.Printf("INF SQL: %s -- %v\n", sql, args)
	return
}
