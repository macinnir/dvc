package query

import (
	"fmt"
	"log"
	"strings"
)

// OrderByDir is a type containing the order by string for ORDER BY sql clauses
type OrderByDir string

// Directional constants
const (
	ASC  OrderByDir = "ASC"
	DESC OrderByDir = "DESC"
)

// Placeholder used for sql arguments
const Placeholder string = "?"

// Select initiates a select statement
// Select(DomainObject)
func Select() *SelectQuery {
	s := &SelectQuery{}
	return s
}

// SelectQuery is the containing struct for functionality to build select statements on domain objects
type SelectQuery struct {
	isDistinct bool
	Object     DomainObject
	where      []IQueryPart
	limit      IQueryPart
	orderBy    []OrderBy
}

// func (s *SelectQuery) Run()

// Distinct sets a flag indicating that the SELECT statement will start with a DISTINCT clause
func (s *SelectQuery) Distinct() *SelectQuery {
	s.isDistinct = true
	return s
}

// Where is sugar for the And(queryParts...) method
func (s *SelectQuery) Where(queryParts ...IQueryPart) *SelectQuery {
	return s.And(queryParts...)
}

// And adds an "and group" to the query
func (s *SelectQuery) And(queryParts ...IQueryPart) *SelectQuery {

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
func (s *SelectQuery) Or(queryParts ...IQueryPart) *SelectQuery {

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
func (s *SelectQuery) Limit(offset int, count int) *SelectQuery {

	s.limit = Limit{offset, count}

	return s
}

// OrderBy adds an ORDER BY clause to the sql statement
func (s *SelectQuery) OrderBy(fieldName string, orderByDir OrderByDir) *SelectQuery {
	if s.orderBy == nil {
		s.orderBy = []OrderBy{}
	}

	s.orderBy = append(s.orderBy, OrderBy{Field: fieldName, Dir: orderByDir})

	return s
}

// ToSQL builds a sql select statement
func (s *SelectQuery) ToSQL() (sql string, args []interface{}) {

	args = []interface{}{}

	fields := s.Object.GetFieldsOrdered()

	escapedFields := []string{}
	for _, field := range fields {
		fmt.Println(field.Name)
		escapedFields = append(escapedFields, escapeField(field.Name))
	}

	wheres := []string{}

	if len(s.where) > 0 {
		for _, where := range s.where {
			whereSQL, whereArgs := where.ToSQL()
			wheres = append(wheres, whereSQL)
			if whereArgs != nil {
				for _, whereArg := range whereArgs {
					args = append(args, whereArg)
				}
			}
		}
	}

	where := ""
	if len(wheres) > 0 {
		where = "WHERE " + strings.Join(wheres, " ") + " "
	}

	if s.orderBy != nil && len(s.orderBy) > 0 {
		where += "ORDER BY "
		orderBys := []string{}
		for _, orderBy := range s.orderBy {
			orderBySQL, _ := orderBy.ToSQL()
			orderBys = append(orderBys, orderBySQL)
		}
		where += strings.Join(orderBys, ", ") + " "
	}

	if s.limit != nil {
		where += "LIMIT "
		limitSQL, limitArgs := s.limit.ToSQL()
		where += limitSQL
		for _, limitArg := range limitArgs {
			args = append(args, limitArg)
		}
	}

	distinct := " "
	if s.isDistinct == true {
		distinct = " DISTINCT "
	}

	sql = fmt.Sprintf("SELECT%s%s FROM %s %s", distinct, strings.Join(escapedFields, ","), escapeField(s.Object.GetName()), where)

	log.Printf("INF SQL: %s -- %v\n", sql, args)
	return
}
