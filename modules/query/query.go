package main

import (
	"fmt"
)

// SELECT * FROM `User` WHERE (
//		DateCreated < NOW() OR IsDeleted != 1
// ) AND LastUpdated > DATE_SUB(NOW(), INTERVAL 1 HOUR)

// Q("Users").Where("UserID").Equals(1);
// Q("Users").Where("DateCreated").LessThan(

type Comparison int

const (
	Equals Comparison = iota
	NotEquals
	GreaterThan
	LesserThan
	GreaterThanOrEquals
	LesserThanOrEquals
	IN
	NOTIN
	AND
	OR
	Limit
	OrderBy
)

// Query is the internal query object
type Query struct {
	tableName    string
	wheres       []*Where
	lastWhereIdx int
}

type Where struct {
	column     string
	column2    string
	value      interface{}
	comparison Comparison
}

func (q *Query) Where(columnName string) *Query {
	q.wheres = append(q.wheres, &Where{column: columnName})
	q.lastWhereIdx = len(q.wheres) - 1
	return q
}

func (q *Query) Equals(value interface{}) *Query {
	if q.lastWhereIdx < 0 {
		panic("Cannot call Equals() unless after a where clause...")
	}

	q.wheres[q.lastWhereIdx].value = value
	q.wheres[q.lastWhereIdx].comparison = Equals
	q.lastWhereIdx = -1
	return q
}

func (q *Query) NotEquals(value interface{}) *Query {
	if q.lastWhereIdx < 0 {
		panic("Cannot call NotEquals() unless after a where clause...")
	}

	q.wheres[q.lastWhereIdx].value = value
	q.wheres[q.lastWhereIdx].comparison = NotEquals
	q.lastWhereIdx = -1
	return q
}

func (q *Query) GreaterThan(value interface{}) *Query {
	if q.lastWhereIdx < 0 {
		panic("Cannot call GreaterThan() unless after a where clause...")
	}

	q.wheres[q.lastWhereIdx].value = value
	q.wheres[q.lastWhereIdx].comparison = GreaterThan
	return q
}
func (q *Query) LesserThan(value interface{}) *Query {
	if q.lastWhereIdx < 0 {
		panic("Cannot call (() unless after a where clause...")
	}

	q.wheres[q.lastWhereIdx].value = value
	q.wheres[q.lastWhereIdx].comparison = LesserThan
	return q
}

func (q *Query) GreaterThanOrEquals(value interface{}) *Query {
	if q.lastWhereIdx < 0 {
		panic("Cannot call GreaterThan() unless after a where clause...")
	}

	q.wheres[q.lastWhereIdx].value = value
	q.wheres[q.lastWhereIdx].comparison = GreaterThan
	return q
}

func (q *Query) LesserThanOrEquals(value interface{}) *Query {
	if q.lastWhereIdx < 0 {
		panic("Cannot call LesserThanO() unless after a where clause...")
	}

	q.wheres[q.lastWhereIdx].value = value
	q.wheres[q.lastWhereIdx].comparison = LesserThanOrEquals
	return q
}

func (q *Query) IN(value interface{}) *Query {
	if q.lastWhereIdx < 0 {
		panic("Cannot call in() unless after a where clause...")
	}

	q.wheres[q.lastWhereIdx].value = value
	q.wheres[q.lastWhereIdx].comparison = IN
	return q
}

func (q *Query) NOTIN(value interface{}) *Query {
	if q.lastWhereIdx < 0 {
		panic("Cannot call value() unless after a where clause...")
	}

	q.wheres[q.lastWhereIdx].value = value
	q.wheres[q.lastWhereIdx].comparison = NOTIN
	return q
}

func (q *Query) AND() *Query {

	if q.lastWhereIdx > -1 {
		panic("Cannot call AND() unless starting a new where clause...")
	}

	q.wheres = append(q.wheres, &Where{
		comparison: AND,
	})

	return q
}

func (q *Query) OR() *Query {
	if q.lastWhereIdx > -1 {
		panic("Cannot call OR() unless starting a new where clause...")
	}

	q.wheres = append(q.wheres, &Where{
		comparison: OR,
	})

	return q
}

func (q *Query) Limit(offset int, count int) *Query {
	if q.lastWhereIdx > -1 {
		panic("Cannot call Limit() unless starting a new where clause...")
	}

	q.wheres = append(q.wheres, &Where{
		comparison: Limit,
		value:      fmt.Sprintf("%d, %d", offset, count),
	})

	return q
}

func (q *Query) OrderBy(field string, direction string) *Query {
	if q.lastWhereIdx > -1 {
		panic("Cannot call Limit() unless starting a new where clause...")
	}

	q.wheres = append(q.wheres, &Where{
		comparison: OrderBy,
		value:      field + " " + direction,
	})

	return q
}

func (q *Query) ToString() string {
	u := "SELECT * FROM `" + q.tableName + "`"

	if len(q.wheres) > 0 {
		u += " WHERE"

		for _, where := range q.wheres {
			switch where.comparison {
			case AND:
				u += " AND"
			case OR:
				u += " OR"
			case Equals:
				u += " `" + where.column + "` = ?"
			case NotEquals:
				u += " `" + where.column + "` != ?"
			case GreaterThan:
				u += " `" + where.column + "` > ?"
			case LesserThan:
				u += " `" + where.column + "` < ?"
			case GreaterThanOrEquals:
				u += " `" + where.column + "` >= ?"
			case LesserThanOrEquals:
				u += " `" + where.column + "` <= ?"
			case Limit:
				u += " LIMIT " + where.value.(string)
			case OrderBy:
				u += " ORDER BY " + where.value.(string)
			}
		}
	}

	return u
}

// Q returns a new Query
func Q(tableName string) *Query {
	return &Query{
		tableName:    tableName,
		wheres:       []*Where{},
		lastWhereIdx: -1,
	}
}

func main() {
	q := Q("foo").
		Where("bar").
		Equals("baz").
		OR().
		Where("quux").
		Equals(1).
		OrderBy("foo", "DESC").
		Limit(1, 5).
		ToString()
	fmt.Println(q)
}
