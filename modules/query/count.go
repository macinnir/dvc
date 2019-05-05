package query

import (
	"fmt"
	"log"
)

// Count initiates a count statement
// Count(DomainObject)
func Count() *CountQuery {
	s := &CountQuery{}
	return s
}

// CountQuery is the containing struct for functionality to build count statements on domain objects
type CountQuery struct {
	Object DomainObject
	where  []IQueryPart
}

// Where is sugar for the And(queryParts...) method
func (s *CountQuery) Where(queryParts ...IQueryPart) *CountQuery {
	return s.And(queryParts...)
}

// And adds an "and group" to the query
func (s *CountQuery) And(queryParts ...IQueryPart) *CountQuery {

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
func (s *CountQuery) Or(queryParts ...IQueryPart) *CountQuery {

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
func (s *CountQuery) ToSQL() (sql string, args []interface{}) {

	sql, args = buildWhereClauseString(s.where)

	sql = fmt.Sprintf("SELECT COUNT(*) FROM %s %s", escapeField(s.Object.GetName()), sql)

	log.Printf("INF SQL: %s -- %v\n", sql, args)
	return
}
