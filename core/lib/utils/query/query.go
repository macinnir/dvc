package query

import (
	"errors"
	"fmt"
	"strings"
)

type Column string
type TableName string

type IModel interface {
	Table_Name() TableName
	Table_Columns() []Column
	Table_PrimaryKey() Column
	Table_PrimaryKey_Value() int64
	Table_InsertColumns() []Column
	Table_UpdateColumns() []Column
	Table_Column_Types() map[Column]string
	String() string
	Update() string
	Create() string
	Destroy() string
	FromID(id int64) string

	// Table_Column_Values() map[string]interface{}
}

type QueryType int

const (
	QueryTypeNotSet QueryType = iota
	QueryTypeSelect
	QueryTypeUpdate
	QueryTypeDelete
	QueryTypeInsert
)

type QueryOrderBy int

const (
	QueryOrderByASC QueryOrderBy = iota
	QueryOrderByDESC
)

func (q QueryOrderBy) String() string {
	switch q {
	case QueryOrderByASC:
		return "ASC"
	default:
		// QueryOrderByDESC
		return "DESC"
	}
}

type Q struct {
	fields      []string
	alias       string
	model       IModel
	queryType   QueryType
	where       *whereClause
	limit       int64
	offset      int64
	orderBy     [][]string
	setSorter   []Column
	sets        map[Column]interface{}
	columnTypes map[Column]string
	errors      []string
	inst        int64
}

func Query(model IModel) *Q {
	return &Q{
		fields:      []string{},
		model:       model,
		orderBy:     [][]string{},
		setSorter:   []Column{},
		sets:        map[Column]interface{}{},
		columnTypes: model.Table_Column_Types(),
		alias:       "t",
		errors:      []string{},
		inst:        0,
	}
}

func (q *Q) Alias(alias string) *Q {
	q.alias = alias
	return q
}

func (q *Q) Limit(limit, offset int64) *Q {
	q.limit = limit
	q.offset = offset
	return q
}

func (q *Q) OrderBy(col string, dir QueryOrderBy) *Q {
	q.orderBy = append(q.orderBy, []string{col, dir.String()})
	return q
}

func (q *Q) Fields(fields ...string) *Q {
	q.fields = fields
	return q
}

// Field includes a specific field in the columns to be returned by a result set
func (q *Q) Field(name Column) *Q {

	if _, ok := q.columnTypes[name]; !ok {
		q.errorInvalidColumn("SELECT", string(name))
	}

	q.fields = append(q.fields, "`"+q.alias+"`.`"+string(name)+"`")

	return q
}

// FieldAs includes a specific field in the columns to be returned by a set aliased by `as`
func (q *Q) FieldAs(name Column, as string) *Q {

	if _, ok := q.columnTypes[name]; !ok {
		q.errorInvalidColumn("SELECT", string(name))
	}

	q.fields = append(q.fields, "`"+q.alias+"`.`"+string(name)+"` AS `"+as+"`")

	return q
}

// FieldRaw allows for an arbitrary string (e.g. "NOW()") to be included in the select columns
func (q *Q) FieldRaw(fieldStr, as string) *Q {
	q.fields = append(q.fields, fieldStr+" AS "+"`"+as+"`")

	return q
}

func (q *Q) Count(name Column, as string) *Q {

	if _, ok := q.columnTypes[name]; !ok {
		q.errorInvalidColumn("COUNT()", string(name))
		return q
	}

	return q.FieldRaw("COUNT(`"+q.alias+"`.`"+string(name)+"`)", as)
}

func (q *Q) Sum(name Column, as string) *Q {

	if _, ok := q.columnTypes[name]; !ok {
		q.errorInvalidColumn("Sum()", string(name))
		return q
	}

	return q.FieldRaw("COALESCE(SUM(`"+q.alias+"`.`"+string(name)+"`), 0)", as)
}

// Where().Equals("a", "b")
func (q *Q) Where(args ...wherePart) *Q {
	// allow for multiple where calls in single query
	if q.where == nil {
		q.where = &whereClause{
			query:      q,
			WhereParts: []wherePart{},
		}
	}

	for k := range args {
		q.where.WhereParts = append(q.where.WhereParts, args[k])
	}
	return q
}

// func Save(model IModel) *Q {

// 	var q *Q
// 	colMap := model.Table_Column_Values()

// 	if model.Table_PrimaryKey_Value() > 0 {
// 		q = Update(model)
// 		updateColumns := model.Table_UpdateColumns()
// 		for _, col := range updateColumns {
// 			q.Set(col, colMap[col])
// 		}
// 		q.Where(
// 			EQ(model.Table_PrimaryKey(), model.Table_PrimaryKey_Value()),
// 		)
// 	} else {
// 		q = Insert(model)
// 		insertColumns := model.Table_InsertColumns()
// 		for _, col := range insertColumns {
// 			q.Set(col, colMap[col])
// 		}
// 	}

// 	return q
// }

// func Destroy(model IModel) *Q {

// 	var q *Q
// 	colMap := model.Table_Column_Values()

// 	if colMap[model.Table_PrimaryKey()].(int64) > 0 {
// 		q = Delete(model)
// 		q.Where(
// 			EQ(model.Table_PrimaryKey(), model.Table_PrimaryKey_Value()),
// 		)
// 	}
// 	return q

// }

func (q *Q) String() (string, error) {

	var sb strings.Builder

	switch q.queryType {
	case QueryTypeSelect:

		fields := "`" + q.alias + "`.*"

		if len(q.fields) > 0 {
			fields = strings.Join(q.fields, ", ")
		}
		sb.WriteString("SELECT " + fields + " FROM")
		// sql += fmt.Sprintf("SEsLECT %s FROM", fields)
	case QueryTypeInsert:
		sb.WriteString("INSERT INTO")
		// sql += "INSERT INTO"
		q.alias = ""
	case QueryTypeUpdate:
		sb.WriteString("UPDATE")
		// sql += "UPDATE"
		q.alias = ""
	case QueryTypeDelete:
		sb.WriteString("DELETE FROM")
		// sql += "DELETE FROM"s
		q.alias = ""
	}

	sb.WriteString(" `" + string(q.model.Table_Name()) + "`")
	// sql += " `" + q.model.Table_Name() + "`"

	if len(q.alias) > 0 && q.queryType == QueryTypeSelect {
		sb.WriteString(" `" + q.alias + "`")
		// sql += " `" + q.alias + "`"
	}

	if q.queryType == QueryTypeUpdate && len(q.sets) > 0 {

		sb.WriteString(" SET ")
		// sql += " SET "

		setStmts := []string{}

		for k := range q.setSorter {

			colName := q.setSorter[k]

			val := fmt.Sprint(q.sets[colName])

			if _, ok := q.columnTypes[colName]; !ok {
				q.errorInvalidColumn("UPDATE()", string(colName))
				val = "'" + val + "'"
			} else {
				if q.columnTypes[colName] == "%s" {
					val = "'" + val + "'"
				}
			}

			setStmts = append(setStmts, q.col(string(colName))+" = "+val)
		}

		sb.WriteString(strings.Join(setStmts, ", "))
		// sql += strings.Join(setStmts, ", ")
	}

	if q.queryType == QueryTypeInsert && len(q.sets) > 0 {

		cols := []string{}
		vals := []string{}

		for k := range q.setSorter {

			colName := q.setSorter[k]

			if _, ok := q.columnTypes[colName]; !ok {
				q.errorInvalidColumn("INSERT()", string(colName))
			}

			val := fmt.Sprint(q.sets[colName])

			if q.columnTypes[colName] == "%s" {
				val = "'" + val + "'"
			}

			cols = append(cols, q.col(string(colName)))
			vals = append(vals, val)
		}

		sb.WriteString(" ( " + strings.Join(cols, ", ") + " ) VALUES ( " + strings.Join(vals, ", ") + " )")
		// sql += " ( " + strings.Join(cols, ", ") + " ) VALUES ( " + strings.Join(vals, ", ") + " )"
	}

	if q.where != nil && q.queryType != QueryTypeInsert {
		sb.WriteString(" WHERE ")
		whereClause := q.printWhereClause(q.columnTypes, q.where.WhereParts)
		if len(whereClause) == 0 {
			fmt.Println(q.where.WhereParts, q.queryType)
			q.error(fmt.Sprintf("EMPTY_WHERE_CLAUSE: `%s`", q.model.Table_Name()))
		}

		sb.WriteString(whereClause)
	}

	if q.queryType == QueryTypeSelect && len(q.orderBy) > 0 {
		orderBys := []string{}
		for k := range q.orderBy {
			orderBys = append(orderBys, q.col(q.orderBy[k][0])+" "+strings.ToUpper(q.orderBy[k][1]))
		}
		sb.WriteString(" ORDER BY " + strings.Join(orderBys, ","))
	}

	if q.limit > 0 {
		sb.WriteString(" LIMIT " + fmt.Sprint(q.limit))
		// sql += fmt.Sprintf(" LIMIT %d", q.limit)
	}

	if q.offset > 0 {
		sb.WriteString(" OFFSET " + fmt.Sprint(q.offset))
		// sql += fmt.Sprintf(" OFFSET %d", q.offset)
	}

	var e error

	if len(q.errors) > 0 {
		e = errors.New(strings.Join(q.errors, "\n--"))
	}

	return sb.String(), e
}

func (q *Q) col(colName string) string {
	if len(q.alias) > 0 {
		return "`" + q.alias + "`.`" + colName + "`"
		// return fmt.Sprintf("`%s`.`%s`", q.alias, colName)
	}
	return "`" + string(colName) + "`"
	// return fmt.Sprintf("`%s`", colName)
}

func isConjunction(whereType WhereType) bool {

	switch whereType {
	case WhereTypeAnd, WhereTypeOr:
		return true
	default:
		return false
	}

}

func (q *Q) printWhereClause(columnTypes map[Column]string, whereParts []wherePart) string {

	sb := strings.Builder{}

	// prevWasConjunction := false

	for k := range whereParts {

		w := whereParts[k]

		if w.e != nil {
			q.error(w.e.Error())
		}

		isConj := isConjunction(w.whereType)

		// If this is is not a conjunction AND fieldName is not empty
		if !isConj && len(w.fieldName) > 0 {

			sb.WriteString(q.col(w.fieldName))

			if _, ok := columnTypes[Column(w.fieldName)]; !ok {
				q.errorInvalidColumn("WHERE()", w.fieldName)
			}
		}

		column := columnTypes[Column(w.fieldName)]

		switch w.whereType {
		case WhereTypeEquals, WhereTypeEqualsField:
			sb.WriteString(" = ")
		case WhereTypeNotEquals:
			sb.WriteString(" <> ")
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
		case WhereTypeExists:
			sb.WriteString("EXISTS")
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
		case WhereTypeAll:
			sb.WriteString("1=1")

		case WhereTypeLike:
			if column != "%s" {
				q.errorInvalidValue("LIKE", column, w.values[0])
			}
			sb.WriteString(" LIKE ")

		case WhereTypeNotLike:
			if column != "%s" {
				q.errorInvalidValue("NOT LIKE", column, w.values[0])
			}
			sb.WriteString(" NOT LIKE ")
		}

		if w.whereType != WhereTypeExists && !isConj && len(w.values) > 0 {

			switch w.whereType {
			case WhereTypeEqualsField:
				sb.WriteString(w.values[0].(string))
			case WhereTypeBetween:
				list := []string{}
				for l := range w.values {
					if column == "%s" {
						list = append(list, "'"+fmt.Sprint(w.values[l])+"'")
					} else {
						list = append(list, fmt.Sprint(w.values[l]))
					}
				}
				sb.WriteString(list[0] + " AND " + list[1])
			case WhereTypeIN:
				list := []string{}
				for l := range w.values {
					if column == "%s" {
						list = append(list, "'"+fmt.Sprint(w.values[l])+"'")
					} else {
						list = append(list, fmt.Sprint(w.values[l]))
					}
				}
				sb.WriteString("( " + strings.Join(list, ", ") + " )")
			default:
				if column == "%s" {
					sb.WriteString("'" + fmt.Sprint(w.values[0]) + "'")
				} else {
					sb.WriteString(fmt.Sprint(w.values[0]))
				}
			}
		}

		if w.whereType == WhereTypeExists {
			sb.WriteString(" ( " + fmt.Sprint(w.values[0]) + " )")
		}

		if len(w.subParts) > 0 {
			sb.WriteString(q.printWhereClause(columnTypes, w.subParts))
		}

	}

	return sb.String()
}

func (q *Q) error(err string) {
	q.errors = append(q.errors, err)
}

func (q *Q) errorInvalidColumn(errType string, name string) {
	q.error(fmt.Sprintf("%s: INVALID COLUMN: `%s`.`%s`", errType, q.model.Table_Name(), name))
}

func (q *Q) errorInvalidValue(errType, name string, value interface{}) {
	q.error(fmt.Sprintf("%s: INVALID VALUE: `%s`.`%s` => %v", errType, q.model.Table_Name(), name, value))
}

type WhereType int

const (
	WhereTypeEquals WhereType = iota
	WhereTypeEqualsField
	WhereTypeNotEquals
	WhereTypeGreaterThan
	WhereTypeLessThan
	WhereTypeGreaterThanOrEqualTo
	WhereTypeLessThanOrEqualTo
	WhereTypeBetween
	WhereTypeLike
	WhereTypeNotLike
	WhereTypeIN
	WhereTypeExists
	WhereTypeAnd
	WhereTypeOr
	WhereTypeParenthesisEnd
	WhereTypeParenthesisStart
	// WhereTypeAll is a WHERE clause of `1=1` used for convenience
	// when conditionally adding WHERE clauses starting with a conjunction (AND/OR,etc)
	// separating them.
	// e.g. SELECT * FROM `Foo` WHERE 1=1
	//      SELECT * FROM `Foo` WHERE 1=1 AND FooID = 123;
	WhereTypeAll
)

type wherePart struct {
	whereType WhereType
	fieldName string
	values    []interface{}
	subParts  []wherePart
	e         error
}

func newWhereParent(whereType WhereType, fieldName string, values []interface{}) wherePart {
	return wherePart{
		whereType: whereType,
		fieldName: fieldName,
		values:    values,
		subParts:  []wherePart{},
	}
}

type whereClause struct {
	query      *Q
	WhereParts []wherePart
}

////
// EXPOSED API
////

// EQ is an equals statement between a table column and a value
func EQ(fieldName Column, value interface{}) wherePart {
	return newWhereParent(
		WhereTypeEquals,
		string(fieldName),
		[]interface{}{
			value,
		},
	)
}

func EQF(fieldName1, fieldName2 string) wherePart {
	return newWhereParent(
		WhereTypeEqualsField,
		fieldName1,
		[]interface{}{fieldName2},
	)
}

func NE(fieldName Column, value interface{}) wherePart {
	return newWhereParent(
		WhereTypeNotEquals,
		string(fieldName),
		[]interface{}{value},
	)
}

func LT(fieldName Column, value interface{}) wherePart {
	return newWhereParent(
		WhereTypeLessThan,
		string(fieldName),
		[]interface{}{value},
	)
}

func GT(fieldName Column, value interface{}) wherePart {
	return newWhereParent(
		WhereTypeGreaterThan,
		string(fieldName),
		[]interface{}{value},
	)
}

func LTOE(fieldName Column, value interface{}) wherePart {
	return newWhereParent(
		WhereTypeLessThanOrEqualTo,
		string(fieldName),
		[]interface{}{value},
	)
}

func GTOE(fieldName Column, value interface{}) wherePart {
	return newWhereParent(
		WhereTypeGreaterThanOrEqualTo,
		string(fieldName),
		[]interface{}{value},
	)
}

// IN is an IN clause
func IN(fieldName Column, values ...interface{}) wherePart {
	return newWhereParent(
		WhereTypeIN,
		string(fieldName),
		values,
	)
}

// INString is a helper function for converting a slice of string arguments into
// a slice of interface arguments, passed into an IN clause and returned
func INString(fieldName Column, values []string) wherePart {
	interfaces := make([]interface{}, len(values))

	for k := range values {
		interfaces[k] = values[k]
	}

	return IN(fieldName, interfaces...)
}

// INInt64 is a helper function for converting a slice of string arguments into
// a slice of interface arguments, passed into an IN clause and returned
func INInt64(fieldName Column, values []int64) wherePart {
	interfaces := make([]interface{}, len(values))

	for k := range values {
		interfaces[k] = values[k]
	}

	return IN(fieldName, interfaces...)
}

func Between(fieldName Column, from, to interface{}) wherePart {
	return newWhereParent(
		WhereTypeBetween,
		string(fieldName),
		[]interface{}{from, to},
	)
}

func Like(fieldName Column, value string) wherePart {
	return newWhereParent(
		WhereTypeLike,
		string(fieldName),
		[]interface{}{value},
	)
}

func NotLike(fieldName Column, value string) wherePart {
	return newWhereParent(
		WhereTypeNotLike,
		string(fieldName),
		[]interface{}{value},
	)
}

func And(args ...wherePart) wherePart {

	and := newWhereParent(WhereTypeAnd, "", []interface{}{})

	if len(args) > 0 {
		and.subParts = append(and.subParts, PS())

		for k := range args {
			and.subParts = append(and.subParts, args[k])
		}

		and.subParts = append(and.subParts, PE())
	}

	return and
}

func Or(args ...wherePart) wherePart {

	or := newWhereParent(WhereTypeOr, "", []interface{}{})

	if len(args) > 0 {
		or.subParts = append(or.subParts, PS())

		for k := range args {
			or.subParts = append(or.subParts, args[k])
		}

		or.subParts = append(or.subParts, PE())
	}

	return or
}

// Parenthesis Start
func PS() wherePart {
	return newWhereParent(
		WhereTypeParenthesisStart,
		"",
		[]interface{}{},
	)
}

// Parenthesis End
func PE() wherePart {
	return newWhereParent(
		WhereTypeParenthesisEnd,
		"",
		[]interface{}{},
	)
}

func WhereAll() wherePart {
	return newWhereParent(
		WhereTypeAll,
		"",
		[]interface{}{},
	)
}

func Exists(clause *Q) wherePart {
	clauseString, e := clause.String()

	w := newWhereParent(
		WhereTypeExists,
		"",
		[]interface{}{clauseString},
	)
	if e != nil {
		w.e = e
	}
	return w
}

func Union(queries ...*Q) (string, error) {

	sqls := []string{}
	for k := range queries {
		query, e := queries[k].String()
		if e != nil {
			return "", e
		}
		sqls = append(sqls, query)
	}

	return strings.Join(sqls, " UNION ALL "), nil
}

func Select(model IModel) *Q {
	q := Query(model)
	q.queryType = QueryTypeSelect
	return q
}

func Update(model IModel) *Q {
	q := Query(model)
	q.queryType = QueryTypeUpdate
	return q
}

func (q *Q) Set(fieldName Column, value interface{}) *Q {
	if _, ok := q.sets[fieldName]; !ok {
		q.sets[fieldName] = value
		q.setSorter = append(q.setSorter, fieldName)
	}
	return q
}

func Delete(model IModel) *Q {
	q := Query(model)
	q.queryType = QueryTypeDelete
	return q
}

func Insert(model IModel) *Q {
	q := Query(model)
	q.queryType = QueryTypeInsert
	return q
}
