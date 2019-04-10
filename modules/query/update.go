package query

import (
	"fmt"
	"log"
)

// Update initiates an update statement
func Update(object DomainObject) *UpdateQuery {
	s := &UpdateQuery{sets: Set{}, Object: object}
	return s
}

// UpdateQuery is the containing struct for functionality to build select statements on domain objects
type UpdateQuery struct {
	Object DomainObject
	sets   Set
	where  []IQueryPart
}

// func (s *UpdateQuery) Run()

// Where is sugar for the And(queryParts...) method
func (s *UpdateQuery) Where(queryParts ...IQueryPart) *UpdateQuery {
	return s.And(queryParts...)
}

// And adds an "and group" to the query
func (s *UpdateQuery) And(queryParts ...IQueryPart) *UpdateQuery {

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
func (s *UpdateQuery) Or(queryParts ...IQueryPart) *UpdateQuery {

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

// Set adds a set part to the fields to be updated
func (s *UpdateQuery) Set(fieldName string, value interface{}) *UpdateQuery {
	s.sets[fieldName] = value
	return s
}

// ToSQL builds a sql select statement
func (s *UpdateQuery) ToSQL() (sql string, args []interface{}) {

	setSQL, setArgs := s.sets.ToSQL()

	args = []interface{}{}
	for _, val := range setArgs {
		args = append(args, val)
	}

	whereSQL, whereArgs := buildWhereClauseString(s.where)

	sql = fmt.Sprintf("UPDATE %s SET %s %s", escapeField(s.Object.GetName()), setSQL, whereSQL)

	for _, val := range whereArgs {
		args = append(args, val)
	}

	log.Printf("INF SQL: %s -- %v\n", sql, args)
	return
}
