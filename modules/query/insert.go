package query

import (
	"fmt"
	"log"
)

// Insert initiates an update statement
func Insert() *InsertQuery {
	s := &InsertQuery{values: Values{}}
	return s
}

// InsertQuery is the containing struct for functionality to build select statements on domain objects
type InsertQuery struct {
	Object DomainObject
	values Values
}

// Value adds a value part for the fields to be inserted
func (s *InsertQuery) Value(fieldName string, value interface{}) *InsertQuery {
	s.values[fieldName] = value
	return s
}

// ToSQL builds a sql select statement
func (s *InsertQuery) ToSQL() (sql string, args []interface{}) {

	var valuesSQL string
	valuesSQL, args = s.values.ToSQL()

	sql = fmt.Sprintf("INSERT INTO %s %s", escapeField(s.Object.GetName()), valuesSQL)

	log.Printf("INF SQL: %s -- %v\n", sql, args)
	return
}
