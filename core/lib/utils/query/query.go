package query

import (
	"errors"
	"fmt"
	"strings"
)

type IModel interface {
	Table_Name() string
	Table_Columns() []string
	Table_PrimaryKey() string
	Table_PrimaryKey_Value() int64
	Table_InsertColumns() []string
	Table_UpdateColumns() []string
	Table_Column_Types() map[string]string
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

type Q struct {
	fields      []string
	alias       string
	model       IModel
	queryType   QueryType
	where       *Where
	limit       int64
	offset      int64
	orderBy     [][]string
	setSorter   []string
	sets        map[string]interface{}
	columnTypes map[string]string
	errors      []string
}

func Query(model IModel) *Q {
	return &Q{
		fields:      []string{},
		model:       model,
		orderBy:     [][]string{},
		setSorter:   []string{},
		sets:        map[string]interface{}{},
		columnTypes: model.Table_Column_Types(),
		alias:       "t",
		errors:      []string{},
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

func (q *Q) OrderBy(col, dir string) *Q {
	q.orderBy = append(q.orderBy, []string{col, dir})
	return q
}

func (q *Q) Fields(fields ...string) *Q {
	q.fields = fields
	return q
}

func (q *Q) Field(name string) *Q {

	if _, ok := q.columnTypes[name]; !ok {
		q.error("Unknown column `" + name + "`")
	}

	q.fields = append(q.fields, "`"+q.alias+"`.`"+name+"`")

	return q
}

func (q *Q) FieldAs(name, as string) *Q {
	q.fields = append(q.fields, "`"+q.alias+"`.`"+name+"` AS `"+as+"`")

	return q
}

func (q *Q) FieldRaw(fieldStr, as string) *Q {
	q.fields = append(q.fields, fieldStr+" AS "+"`"+as+"`")

	return q
}

func (q *Q) Count(name, as string) *Q {

	if _, ok := q.columnTypes[name]; !ok {
		q.errors = append(q.errors, fmt.Sprintf("COUNT(): Invalid column `%s`", name))
		return q
	}

	return q.FieldRaw("COUNT(`"+q.alias+"`.`"+name+"`)", as)
}

func (q *Q) Sum(name, as string) *Q {

	if _, ok := q.columnTypes[name]; !ok {
		q.errors = append(q.errors, fmt.Sprintf("SUM(): Invalid column `%s`", name))
		return q
	}

	return q.FieldRaw("COALESCE(SUM(`"+q.alias+"`.`"+name+"`), 0)", as)
}

// Where().Equals("a", "b")
func (q *Q) Where(args ...WherePart) *Q {
	q.where = &Where{
		query:      q,
		WhereParts: []WherePart{},
	}

	for k := range args {
		q.where.WhereParts = append(q.where.WhereParts, args[k])
	}
	return q
}

func (q *Q) error(err string) {
	q.errors = append(q.errors, fmt.Sprintf("SQL ERROR: Table `%s`: %s", q.model.Table_Name(), err))
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

	sb.WriteString(" `" + q.model.Table_Name() + "`")
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
				q.error(fmt.Sprintf("UPDATE(): Invalid column `%s`", colName))
				val = "'" + val + "'"
			} else {
				if q.columnTypes[colName] == "%s" {
					val = "'" + val + "'"
				}
			}

			setStmts = append(setStmts, q.col(colName)+" = "+val)
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
				q.error(fmt.Sprintf("INSERT(): Invalid column `%s`", colName))
			}

			val := fmt.Sprint(q.sets[colName])

			if q.columnTypes[colName] == "%s" {
				val = "'" + val + "'"
			}

			cols = append(cols, q.col(colName))
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
			q.error(fmt.Sprintf("Empty where clause: %s", whereClause))
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
	return "`" + colName + "`"
	// return fmt.Sprintf("`%s`", colName)
}

func (q *Q) printWhereClause(columnTypes map[string]string, whereParts []WherePart) string {

	sb := strings.Builder{}

	for k := range whereParts {

		w := whereParts[k]

		if w.e != nil {
			q.error(w.e.Error())
		}

		if !(w.whereType == WhereTypeAnd || w.whereType == WhereTypeOr || len(w.fieldName) == 0) {

			sb.WriteString(q.col(w.fieldName))

			if _, ok := columnTypes[w.fieldName]; !ok {
				q.error("Invalid field name `" + w.fieldName + "`")
			}
		}

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
		}

		column := columnTypes[w.fieldName]

		if !(w.whereType == WhereTypeExists ||
			w.whereType == WhereTypeAnd ||
			w.whereType == WhereTypeOr ||
			len(w.values) == 0) {

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

type WherePart struct {
	whereType WhereType
	fieldName string
	values    []interface{}
	subParts  []WherePart
	e         error
}

func NewWherePart(whereType WhereType, fieldName string, values []interface{}) WherePart {
	return WherePart{
		whereType: whereType,
		fieldName: fieldName,
		values:    values,
		subParts:  []WherePart{},
	}
}

type Where struct {
	query      *Q
	WhereParts []WherePart
}

func EQ(fieldName string, value interface{}) WherePart {
	return NewWherePart(
		WhereTypeEquals,
		fieldName,
		[]interface{}{
			value,
		},
	)
}

func EQF(fieldName1, fieldName2 string) WherePart {
	return NewWherePart(
		WhereTypeEqualsField,
		fieldName1,
		[]interface{}{fieldName2},
	)
}

func NE(fieldName string, value interface{}) WherePart {
	return NewWherePart(
		WhereTypeNotEquals,
		fieldName,
		[]interface{}{value},
	)
}

func LT(fieldName string, value interface{}) WherePart {
	return NewWherePart(
		WhereTypeLessThan,
		fieldName,
		[]interface{}{value},
	)
}

func GT(fieldName string, value interface{}) WherePart {
	return NewWherePart(
		WhereTypeGreaterThan,
		fieldName,
		[]interface{}{value},
	)
}

func LTOE(fieldName string, value interface{}) WherePart {
	return NewWherePart(
		WhereTypeLessThanOrEqualTo,
		fieldName,
		[]interface{}{value},
	)
}

func GTOE(fieldName string, value interface{}) WherePart {
	return NewWherePart(
		WhereTypeGreaterThanOrEqualTo,
		fieldName,
		[]interface{}{value},
	)
}

func IN(fieldName string, values ...interface{}) WherePart {
	return NewWherePart(
		WhereTypeIN,
		fieldName,
		values,
	)
}

func Between(fieldName string, from, to interface{}) WherePart {
	return NewWherePart(
		WhereTypeBetween,
		fieldName,
		[]interface{}{from, to},
	)
}

func And(args ...WherePart) WherePart {

	and := NewWherePart(WhereTypeAnd, "", []interface{}{})

	if len(args) > 0 {
		and.subParts = append(and.subParts, PS())

		for k := range args {
			and.subParts = append(and.subParts, args[k])
		}

		and.subParts = append(and.subParts, PE())
	}

	return and
}

func Or(args ...WherePart) WherePart {

	or := NewWherePart(WhereTypeOr, "", []interface{}{})

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
func PS() WherePart {
	return NewWherePart(
		WhereTypeParenthesisStart,
		"",
		[]interface{}{},
	)
}

// Parenthesis End
func PE() WherePart {
	return NewWherePart(
		WhereTypeParenthesisEnd,
		"",
		[]interface{}{},
	)
}

func WhereAll() WherePart {
	return NewWherePart(
		WhereTypeAll,
		"",
		[]interface{}{},
	)
}

func Exists(clause *Q) WherePart {
	clauseString, e := clause.String()

	w := NewWherePart(
		WhereTypeExists,
		"",
		[]interface{}{clauseString},
	)
	if e != nil {
		w.e = e
	}
	return w
}

func CommentID(args ...interface{}) WherePart {
	if len(args) > 1 {
		return IN("CommentID", args...)
	} else {
		return EQ("CommentID", args[0])
	}
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

func (q *Q) Set(fieldName string, value interface{}) *Q {
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
