package query

import (
	"fmt"
	"log"
)

// Delete initiates a select statement
func Delete() *DeleteQuery {
	s := &DeleteQuery{}
	return s
}

// DeleteQuery is the containing struct for functionality to build select statements on domain objects
type DeleteQuery struct {
	Object DomainObject
	where  []IQueryPart
	limit  IQueryPart
}

// Where is sugar for the And(queryParts...) method
func (s *DeleteQuery) Where(queryParts ...IQueryPart) *DeleteQuery {
	return s.And(queryParts...)
}

// And adds an "and group" to the query
func (s *DeleteQuery) And(queryParts ...IQueryPart) *DeleteQuery {

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
func (s *DeleteQuery) Or(queryParts ...IQueryPart) *DeleteQuery {

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

// Limit adds a limit clause to the sql statement
func (s *DeleteQuery) Limit(count int) *DeleteQuery {

	s.limit = LimitSimple(count)

	return s
}

// ToSQL builds a sql select statement
func (s *DeleteQuery) ToSQL() (sql string, args []interface{}) {

	fields := s.Object.GetFieldsOrdered()

	escapedFields := []string{}
	for _, field := range fields {
		escapedFields = append(escapedFields, escapeField(field.Name))
	}

	sql, args = buildWhereClauseString(s.where)

	if s.limit != nil {
		sql += "LIMIT "
		limitSQL, limitArgs := s.limit.ToSQL()
		sql += limitSQL
		args = append(args, limitArgs[0])
	}

	sql = fmt.Sprintf("DELETE FROM %s %s", escapeField(s.Object.GetName()), sql)

	log.Printf("INF SQL: %s -- %v\n", sql, args)
	return
}
