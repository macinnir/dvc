package query

import (
	"fmt"
	"strings"
)

// Equals represents an equals comparison query part
type Equals map[string]interface{}

// ToSQL returns the sql representation of an equals query part along with its arguments (if any)
func (equals Equals) ToSQL() (sql string, args []interface{}) {
	keys := getSortedKeys(equals)
	args = []interface{}{}

	sqls := []string{}
	for _, key := range keys {
		val := equals[key]
		if val == nil {
			sqls = append(sqls, fmt.Sprintf("%s IS NULL", escapeField(key)))
		} else {
			sqls = append(sqls, fmt.Sprintf("%s = %s", escapeField(key), Placeholder))
			args = append(args, val)
		}
	}

	sql = buildWhereClause(sqls)
	return
}

// NotEquals represents a not equals comparison query part
type NotEquals map[string]interface{}

// ToSQL returns the sql representation of a not equals query part along with its arguments (if any)
func (equals NotEquals) ToSQL() (sql string, args []interface{}) {
	keys := getSortedKeys(equals)
	args = []interface{}{}

	sqls := []string{}

	for _, key := range keys {
		val := equals[key]
		if val == nil {
			sqls = append(sqls, fmt.Sprintf("%s IS NOT NULL", escapeField(key)))
		} else {
			sqls = append(sqls, fmt.Sprintf("%s != %s", escapeField(key), Placeholder))
			args = append(args, val)
		}
	}

	sql = buildWhereClause(sqls)
	return
}

// LessThan represents a less than comparison query part
type LessThan map[string]interface{}

// ToSQL returns the sql representation of a less than query part along with its arguments (if any)
func (lessThan LessThan) ToSQL() (sql string, args []interface{}) {
	keys := getSortedKeys(lessThan)
	args = []interface{}{}

	sqls := []string{}
	for _, key := range keys {
		val := lessThan[key]
		sqls = append(sqls, fmt.Sprintf("%s < %s", escapeField(key), Placeholder))
		args = append(args, val)
	}

	sql = buildWhereClause(sqls)
	return
}

// LessThanOrEqual represents a less than or equal comparison query part
type LessThanOrEqual map[string]interface{}

// ToSQL returns the sql representation of a less than or equal query part along with its arguments (if any)
func (lessThanOrEqual LessThanOrEqual) ToSQL() (sql string, args []interface{}) {
	keys := getSortedKeys(lessThanOrEqual)
	args = []interface{}{}

	sqls := []string{}
	for _, key := range keys {
		val := lessThanOrEqual[key]
		sqls = append(sqls, fmt.Sprintf("%s <= %s", escapeField(key), Placeholder))
		args = append(args, val)
	}

	sql = buildWhereClause(sqls)
	return
}

// GreaterThan represents a greater than comparison query part
type GreaterThan map[string]interface{}

// ToSQL returns the sql representation of a greater than query part along with its arguments (if any)
func (greaterThan GreaterThan) ToSQL() (sql string, args []interface{}) {
	keys := getSortedKeys(greaterThan)
	args = []interface{}{}

	sqls := []string{}
	for _, key := range keys {
		val := greaterThan[key]
		sqls = append(sqls, fmt.Sprintf("%s > %s", escapeField(key), Placeholder))
		args = append(args, val)
	}

	sql = buildWhereClause(sqls)
	return
}

// GreaterThanOrEqual represents a greater than or equal comparison query part
type GreaterThanOrEqual map[string]interface{}

// ToSQL returns the sql representation of a greater than or equal query part along with its arguments (if any)
func (greaterThanOrEqual GreaterThanOrEqual) ToSQL() (sql string, args []interface{}) {
	keys := getSortedKeys(greaterThanOrEqual)
	args = []interface{}{}

	sqls := []string{}
	for _, key := range keys {
		val := greaterThanOrEqual[key]
		sqls = append(sqls, fmt.Sprintf("%s >= %s", escapeField(key), Placeholder))
		args = append(args, val)
	}

	sql = buildWhereClause(sqls)
	return
}

// In represents an IN comparison query part
type In map[string][]interface{}

// ToSQL returns the sql representation of an IN query part along with its arguments (if any)
func (in In) ToSQL() (sql string, args []interface{}) {
	keys := getSortedInKeys(in)
	args = []interface{}{}
	sqls := []string{}
	for _, key := range keys {
		vals := in[key]
		sqls = append(sqls, fmt.Sprintf("%s IN (%s)", escapeField(key), strings.Repeat(","+Placeholder, len(vals))[1:]))
		for _, v := range vals {
			args = append(args, v)
		}
	}
	sql = buildWhereClause(sqls)
	return
}

// NotIn represents a NOT IN comparison query part
type NotIn map[string][]interface{}

// ToSQL returns the sql representation of a NOT IN query part along with its arguments (if any)
func (notIn NotIn) ToSQL() (sql string, args []interface{}) {
	keys := getSortedInKeys(notIn)
	args = []interface{}{}
	sqls := []string{}
	for _, key := range keys {
		vals := notIn[key]
		sqls = append(sqls, fmt.Sprintf("%s NOT IN (%s)", escapeField(key), strings.Repeat(","+Placeholder, len(vals))[1:]))
		for _, v := range vals {
			args = append(args, v)
		}
	}
	sql = buildWhereClause(sqls)
	return
}

// Between represents a BETWEEN query part
type Between map[string][]interface{}

// ToSQL returns the sql representation of a BETWEEN query part along with its arguments (if any)
func (between Between) ToSQL() (sql string, args []interface{}) {

	keys := getSortedInKeys(between)
	args = []interface{}{}
	sqls := []string{}
	for _, key := range keys {
		vals := between[key]
		sqls = append(sqls, fmt.Sprintf("%s BETWEEN %s AND %s", escapeField(key), Placeholder, Placeholder))
		n := 0
		for n < 2 {
			args = append(args, vals[n])
			n++
		}
	}
	sql = buildWhereClause(sqls)
	return
}

// NotBetween represents a NOT BETWEEN query part
type NotBetween map[string][]interface{}

// ToSQL returns the sql representation of a NOT BETWEEN query part along with its arguments (if any)
func (notBetween NotBetween) ToSQL() (sql string, args []interface{}) {

	keys := getSortedInKeys(notBetween)
	args = []interface{}{}
	sqls := []string{}
	for _, key := range keys {
		vals := notBetween[key]
		sqls = append(sqls, fmt.Sprintf("%s NOT BETWEEN %s AND %s", escapeField(key), Placeholder, Placeholder))
		n := 0
		for n < 2 {
			args = append(args, vals[n])
			n++
		}
	}
	sql = buildWhereClause(sqls)
	return
}

// And represents an AND comparison query part
type And struct{}

// ToSQL returns the sql representation of an AND query part along with its arguments (if any)
func (and And) ToSQL() (sql string, args []interface{}) {
	sql = "AND"
	// args = []interface{}{}
	return
}

// Or represents an OR comparison query part
type Or struct{}

// ToSQL returns the sql representation of an OR query part along with its arguments (if any)
func (or Or) ToSQL() (sql string, args []interface{}) {
	sql = "OR"
	return
}

// OpenParenthesis represents an open parenthesis (group start) query part
type OpenParenthesis struct{}

// ToSQL returns the sql representation of an Open Parenthesis (group start) query part
func (openParenthesis OpenParenthesis) ToSQL() (sql string, args []interface{}) {
	sql = "("
	return
}

// CloseParenthesis represents a close parenthesis (group end) query part
type CloseParenthesis struct{}

// ToSQL returns the sql representation of an Close Parenthesis (group end) query part
func (closeParenthesis CloseParenthesis) ToSQL() (sql string, args []interface{}) {
	sql = ")"
	return
}

// Distinct represents a DISTINCT clause query part
type Distinct struct{}

// ToSQL returns the sql representation of a DISTINCT clause
func (distinct Distinct) ToSQL() (sql string, args []interface{}) {
	sql = "DISTINCT"
	return
}

// Limit represents a LIMIT clause query part
type Limit []int

// ToSQL returns the sql representation of a LIMIT clause, along with its arguments
func (limit Limit) ToSQL() (sql string, args []interface{}) {

	l := len(limit)
	if l == 0 || l > 2 {
		panic(fmt.Sprintf("Invalid number of arguments (%d) for limit clause", l))
	}

	sqls := []string{}
	args = []interface{}{}
	for _, i := range limit {
		args = append(args, i)
		sqls = append(sqls, Placeholder)
	}

	sql = strings.Join(sqls, ",")

	return
}

// LimitSimple is a single integer limit
type LimitSimple int

// ToSQL returns the SQL representation of a LIMIT clause, along with its arguments
func (dl LimitSimple) ToSQL() (sql string, args []interface{}) {
	sql = Placeholder
	args = []interface{}{dl}
	return
}

// OrderBy represents an ORDER BY sql clause
type OrderBy struct {
	Field string
	Dir   OrderByDir
}

// ToSQL returns the sql representation of an ORDER BY clause
func (orderBy OrderBy) ToSQL() (sql string, args []interface{}) {

	sql = fmt.Sprintf("%s %s", escapeField(orderBy.Field), string(orderBy.Dir))
	return
}

// Set represents a SQL set statement
type Set map[string]interface{}

// ToSQL returns the SQL representation of a set statement (for update)
func (set Set) ToSQL() (sql string, args []interface{}) {
	args = []interface{}{}
	sqls := []string{}

	keys := getSortedKeys(set)
	for _, key := range keys {
		sqls = append(sqls, fmt.Sprintf("%s = %s", escapeField(key), Placeholder))
		args = append(args, set[key])
	}

	sql = strings.Join(sqls, ", ")
	return
}

// Values represents the SQL values in an insert statement
type Values map[string]interface{}

// ToSQL returns the SQL representation of values in an insert statement
func (values Values) ToSQL() (sql string, args []interface{}) {
	args = []interface{}{}
	fields := []string{}
	placeholders := []string{}
	keys := getSortedKeys(values)
	for _, key := range keys {
		fields = append(fields, fmt.Sprintf("%s", escapeField(key)))
		placeholders = append(placeholders, Placeholder)
		args = append(args, values[key])
	}

	sql = fmt.Sprintf("(%s) VALUES (%s)", strings.Join(fields, ", "), strings.Join(placeholders, ", "))
	return
}
