package query

import (
	"fmt"
	"strings"
)

type WhereType int

const (
	WhereTypeEquals WhereType = iota
	WhereTypeEqualsField
	WhereTypeNotEquals
	WhereTypeNotEqualsField
	WhereTypeGreaterThan
	WhereTypeLessThan
	WhereTypeGreaterThanOrEqualTo
	WhereTypeLessThanOrEqualTo
	WhereTypeBetween
	WhereTypeLike
	WhereTypeNotLike
	WhereTypeIsNull
	WhereTypeIsNotNull
	WhereTypeIN
	WhereTypeNotIN
	WhereTypeExists
	WhereTypeNotExists
	WhereTypeAnd
	WhereTypeOr
	WhereTypeParenthesisEnd
	WhereTypeParenthesisStart
	// WhereTypeNone indicates that the wherePart is a noop for the query,
	// If, however, it contains any child clauses, they will be parsed as individual wherePart objects
	WhereTypeNone
	// WhereTypeAll is a WHERE clause of `1=1` used for convenience
	// when conditionally adding WHERE clauses starting with a conjunction (AND/OR,etc)
	// separating them.
	// e.g. SELECT * FROM `Foo` WHERE 1=1
	//      SELECT * FROM `Foo` WHERE 1=1 AND FooID = 123;
	WhereTypeAll
	WhereTypeMod
	WhereTypeModF
	WhereTypeBitAnd
	WhereTypeRaw
)

// WherePart is a part of a where clause.
// This object is an exposed part of the api to make conditional queries easier
// EXAMPLE:
//
//	wheres := []query.WherePart{
//		query.EQ(models.ObjectRelationship_Column_IsDeleted, 0),
//	}
//	if objectTypeFrom != constants.ObjectTypeUnknown {
//		wheres = append(wheres, query.And(), query.EQ(models.ObjectRelationship_Column_ObjectTypeFrom, objectTypeFrom))
//	}
//	if objectIDFrom > 0 {
//		wheres = append(wheres, query.And(), query.EQ(models.ObjectRelationship_Column_ObjectIDFrom, objectIDFrom))
//	}
type WherePart struct {
	whereType WhereType
	fieldName string
	values    []interface{}
	subParts  []*WherePart
	e         error
}

func newWherePart(whereType WhereType, fieldName string, values []interface{}) *WherePart {
	return &WherePart{
		whereType: whereType,
		fieldName: fieldName,
		values:    values,
		subParts:  []*WherePart{},
	}
}

type whereClause struct {
	query      *Q
	WhereParts []*WherePart
}

// And is an and statement with optional args that, if provided, are wrapped in parentheses
// Example: And() will result in the word `AND` being added to the where clause
// Example: And(EQ(1, 1), And(), And(2, 2)) will result in `AND ( 1 = 1 AND 2 = 2 )`
func And(args ...*WherePart) *WherePart {

	and := newWherePart(WhereTypeAnd, "", []interface{}{})

	if len(args) > 0 {
		and.subParts = append(and.subParts, PS())

		for k := range args {
			and.subParts = append(and.subParts, args[k])
		}

		and.subParts = append(and.subParts, PE())
	}

	return and
}

// Ands takes a list of args and separes them all by `AND`
// Example: Ands(query.EQ(1,1), query.EQ(2,2), query.EQ(3,3)) == 1 = 1 AND 2 = 2 AND 3 = 3
func Ands(args ...*WherePart) *WherePart {

	if len(args) == 0 {
		return nil
	}

	if len(args) == 1 {
		return args[0]
	}

	ands := newWherePart(WhereTypeNone, "", []interface{}{})

	subParts := []*WherePart{}

	for k := range args {

		if args[k] == nil {
			continue
		}

		subParts = append(subParts, args[k])
	}

	for k := range subParts {

		ands.subParts = append(ands.subParts, subParts[k])

		// Last item
		if k == len(subParts)-1 {
			break
		}

		ands.subParts = append(ands.subParts, And())

	}

	return ands
}

func Or(args ...*WherePart) *WherePart {

	or := newWherePart(WhereTypeOr, "", []interface{}{})

	if len(args) > 0 {
		or.subParts = append(or.subParts, PS())

		for k := range args {
			or.subParts = append(or.subParts, args[k])
		}

		or.subParts = append(or.subParts, PE())
	}

	return or
}

// Ors takes a list of args and separes them all by `OR`
// Example: Ors(query.EQ(1,1), query.EQ(2,2), query.EQ(3,3)) == 1 = 1 OR 2 = 2 OR 3 = 3
func Ors(args ...*WherePart) *WherePart {

	if len(args) == 0 {
		return nil
	}

	if len(args) == 1 {
		return args[0]
	}

	ors := newWherePart(WhereTypeNone, "", []interface{}{})

	subParts := []*WherePart{}

	for k := range args {

		if args[k] == nil {
			continue
		}

		subParts = append(subParts, args[k])
	}

	for k := range subParts {

		ors.subParts = append(ors.subParts, subParts[k])

		// Last item
		if k == len(subParts)-1 {
			break
		}

		ors.subParts = append(ors.subParts, Or())
	}

	return ors
}

// Paren adds parenthesis to a query where clause
// .Paren(a, b, c) => (a, b, c)
func Paren(args ...*WherePart) *WherePart {
	n := newWherePart(WhereTypeNone, "", []interface{}{})

	if len(args) > 0 {
		n.subParts = append(n.subParts, PS())
		for k := range args {
			n.subParts = append(n.subParts, args[k])
		}

		n.subParts = append(n.subParts, PE())
	}

	return n
}

// Parenthesis Start
func PS() *WherePart {
	return newWherePart(
		WhereTypeParenthesisStart,
		"",
		[]interface{}{},
	)
}

// Parenthesis End
func PE() *WherePart {
	return newWherePart(
		WhereTypeParenthesisEnd,
		"",
		[]interface{}{},
	)
}

// WhereAll adds a WHERE clause of `1=1` used for convenience
// when conditionally adding WHERE clauses starting with a conjunction (AND/OR,etc)
// separating them.
// e.g. SELECT * FROM `Foo` WHERE 1=1
//
//	SELECT * FROM `Foo` WHERE 1=1 AND FooID = 123;
func WhereAll() *WherePart {
	return newWherePart(
		WhereTypeAll,
		"",
		[]interface{}{},
	)
}

// Where creates or adds to an existing where clause
//
//   - Simple
//     q.Where(query.EQ(query.Column("A"), "B"))
//     WHERE `t`.`B` = `t`.`B`
//
//   - Multiple Arguments
//     q.Where(query.EQ(query.Column("A"), "B"), query.And(), query.EQ(query.Column("C"), "D"))
//     WHERE `t`.`A` = 'B' AND `t`.`C` = 'D'
//
//   - Daisy Chain
//     q.Where(query.EQ(query.Column("A"), "B")).Where(query.And()).Where(query.EQ(query.Column("C"), "D"))
//     WHERE `t`.`A` = 'B' AND `t`.`C` = 'D'
//
//   - Separate lines
//     q.Where(query.EQ(query.Column("A"), "B"))
//     q.Where(query.And())
//     q.Where(query.EQ(query.Column("C"), "D"))
//     WHERE `t`.`A` = 'B' AND `t`.`C` = 'D'
func (q *Q) Where(args ...*WherePart) *Q {
	// allow for multiple where calls in single query
	if q.where == nil {
		q.where = &whereClause{
			query:      q,
			WhereParts: []*WherePart{},
		}
	}

	for k := range args {
		q.where.WhereParts = append(q.where.WhereParts, args[k])
	}
	return q
}

func (q *Q) printWhereClause(columnTypes map[Column]string, whereParts []*WherePart) string {

	sb := strings.Builder{}

	for k := range whereParts {

		if whereParts[k] == nil {
			continue
		}

		w := whereParts[k]

		if w.e != nil {
			q.error(w.e.Error())
		}

		isConj := isConjunction(w.whereType)

		// If this is is not a conjunction AND fieldName is not empty
		if !isConj && len(w.fieldName) > 0 {

			if w.whereType != WhereTypeMod &&
				w.whereType != WhereTypeModF &&
				w.whereType != WhereTypeBitAnd &&
				w.whereType != WhereTypeRaw {
				sb.WriteString(q.col(w.fieldName))
			}

			if _, ok := columnTypes[Column(w.fieldName)]; !ok {
				q.errorInvalidColumn(QUERY_ERROR_INVALID_COLUMN, "WHERE...", w.fieldName)
			}
		}

		column := columnTypes[Column(w.fieldName)]

		switch w.whereType {
		case WhereTypeEquals, WhereTypeEqualsField:
			sb.WriteString(" = ")
		case WhereTypeNotEquals, WhereTypeNotEqualsField:
			sb.WriteString(" != ")
		case WhereTypeGreaterThan:
			sb.WriteString(" > ")
		case WhereTypeLessThan:
			sb.WriteString(" < ")
		case WhereTypeGreaterThanOrEqualTo:
			sb.WriteString(" >= ")
		case WhereTypeLessThanOrEqualTo:
			sb.WriteString(" <= ")
		case WhereTypeIN:
			sb.WriteString(" IN ")
		case WhereTypeNotIN:
			sb.WriteString(" NOT IN ")
		case WhereTypeExists:
			sb.WriteString("EXISTS")
		case WhereTypeNotExists:
			sb.WriteString("NOT EXISTS")
		case WhereTypeBetween:
			sb.WriteString(" BETWEEN ")
		case WhereTypeAnd:
			sb.WriteString(" AND ")
		case WhereTypeOr:
			sb.WriteString(" OR ")
		case WhereTypeParenthesisEnd:
			sb.WriteString(" )")
		case WhereTypeParenthesisStart:
			sb.WriteString("( ")
		case WhereTypeNone:
		case WhereTypeAll:
			sb.WriteString("1=1")

		case WhereTypeLike:
			if column != "%s" {
				q.errorInvalidColumn(QUERY_ERROR_INVALID_VALUE, "WHERE...LIKE", "`"+column+"` value: "+fmt.Sprint(w.values[0]))
			}
			sb.WriteString(" LIKE ")

		case WhereTypeNotLike:
			if column != "%s" {
				q.errorInvalidColumn(QUERY_ERROR_INVALID_VALUE, "WHERE...NOT LIKE", "`"+column+"` value: "+fmt.Sprint(w.values[0]))
			}
			sb.WriteString(" NOT LIKE ")
		case WhereTypeIsNull:
			sb.WriteString(" IS NULL")
		case WhereTypeIsNotNull:
			sb.WriteString(" IS NOT NULL")
		}

		if w.whereType != WhereTypeExists && w.whereType != WhereTypeNotExists && !isConj && len(w.values) > 0 {

			switch w.whereType {
			case WhereTypeEqualsField, WhereTypeNotEqualsField:
				sb.WriteString(w.values[0].(string))
			case WhereTypeMod:
				sb.WriteString(
					"MOD(" + string(q.col(w.fieldName)) + ", " + fmt.Sprint(w.values[0]) + ") = " + fmt.Sprint(w.values[1]),
				)
			case WhereTypeModF:
				sb.WriteString(
					"MOD(" + fmt.Sprint(w.values[0]) + ", " + string(q.col(w.fieldName)) + ") = " + fmt.Sprint(w.values[1]),
				)
			case WhereTypeBitAnd:
				sb.WriteString(
					string(q.col(w.fieldName)) + " & " + fmt.Sprint(w.values[0]) + " = " + fmt.Sprint(w.values[1]),
				)
			case WhereTypeBetween:
				list := []string{}
				for l := range w.values {
					// String
					if column == "%s" {
						list = append(list, "'"+EscapeString(fmt.Sprint(w.values[l]))+"'")
					} else {
						list = append(list, fmt.Sprint(w.values[l]))
					}
				}
				sb.WriteString(list[0] + " AND " + list[1])
			case WhereTypeIN, WhereTypeNotIN:
				list := []string{}
				for l := range w.values {
					// String
					if column == "%s" {
						list = append(list, "'"+EscapeString(fmt.Sprint(w.values[l]))+"'")
					} else {
						list = append(list, fmt.Sprint(w.values[l]))
					}
				}
				sb.WriteString("( " + strings.Join(list, ", ") + " )")
			case WhereTypeRaw:
				sb.WriteString(fmt.Sprint(w.values[0]))
			default:
				// String
				if column == "%s" {
					sb.WriteString("'" + EscapeString(fmt.Sprint(w.values[0])) + "'")
				} else {
					sb.WriteString(fmt.Sprint(w.values[0]))
				}
			}
		}

		if w.whereType == WhereTypeExists || w.whereType == WhereTypeNotExists {
			sb.WriteString(" ( " + fmt.Sprint(w.values[0]) + " )")
		}

		if len(w.subParts) > 0 {
			sb.WriteString(q.printWhereClause(columnTypes, w.subParts))
		}

	}

	return sb.String()
}

func isConjunction(whereType WhereType) bool {

	switch whereType {
	case WhereTypeAnd, WhereTypeOr:
		return true
	default:
		return false
	}

}
