package query

import (
	"errors"
	"fmt"
	"strings"

	"github.com/macinnir/dvc/core/lib/utils/db"
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
	Update(db db.IDB) error
	Create(db db.IDB) error
	Delete(db db.IDB) error
	FromID(db db.IDB, id int64) (IModel, error)

	// Table_Column_Values() map[string]interface{}
}

type QueryType int

const (
	QueryTypeNotSet QueryType = iota
	QueryTypeSelect
	QueryTypeRaw
	QueryTypeUpdate
	QueryTypeDelete
	QueryTypeInsert
)

type OrderDir int

const (
	OrderDirASC OrderDir = iota
	OrderDirDESC
)

func OrderDirFromString(s string) OrderDir {
	s = strings.ToLower(s)
	if s == "desc" {
		return OrderDirDESC
	}

	return OrderDirASC
}

func (q OrderDir) String() string {
	switch q {
	case OrderDirASC:
		return "ASC"
	default:
		// OrderDirDESC
		return "DESC"
	}
}

type FieldType int

const (
	FieldTypeBasic FieldType = iota
	FieldTypeRaw
	FieldTypeCount
	FieldTypeSum
	FieldTypeAvg
	FieldTypeMin
	FieldTypeMax
)

type Field struct {
	FieldType FieldType
	Name      Column
	As        string
	Raw       string
}

// NewField creates a new field.
// 	NewField(FieldTypeBasic, "Foo")
// 	NewField(FieldTypeBasic, "Foo", "Bar") <-- `Foo` AS `Bar`
func NewField(fieldType FieldType, column Column, opts ...string) *Field {

	as := ""

	if len(opts) > 0 {
		as = opts[0]
	}

	return &Field{
		FieldType: fieldType,
		Name:      column,
		As:        as,
		Raw:       "",
	}
}

// NewRawField creates a new field.
// 	NewRawField("`t`.`Foo` AS `Bar`)
func NewRawField(raw string) *Field {

	return &Field{
		FieldType: FieldTypeRaw,
		Name:      "",
		As:        "",
		Raw:       raw,
	}
}

type Q struct {
	fields      []*Field
	noAlias     []int
	alias       string
	raw         string
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
		fields:      []*Field{},
		noAlias:     []int{},
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

func (q *Q) FromFieldToString(field *Field) string {

	as := ""

	if len(field.As) > 0 {
		as = " AS `" + field.As + "`"
	}

	switch field.FieldType {
	case FieldTypeCount:
		return "COUNT(`" + q.alias + "`.`" + string(field.Name) + "`)" + as
	case FieldTypeSum:
		return "COALESCE(SUM(`" + q.alias + "`.`" + string(field.Name) + "`), 0)" + as
	case FieldTypeAvg:
		return "COALESCE(AVG(`" + q.alias + "`.`" + string(field.Name) + "`), 0)" + as
	case FieldTypeMin:
		return "COALESCE(MIN(`" + q.alias + "`.`" + string(field.Name) + "`), 0)" + as
	case FieldTypeMax:
		return "COALESCE(MAX(`" + q.alias + "`.`" + string(field.Name) + "`), 0)" + as
	case FieldTypeRaw:
		return field.Raw
	// FieldTypeBasic
	default:
		return "`" + q.alias + "`.`" + string(field.Name) + "`" + as
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

func (q *Q) LimitPage(limit, page int64) *Q {
	q.limit = limit
	q.offset = limit * page
	return q
}

func (q *Q) OrderBy(col Column, dir OrderDir) *Q {
	q.orderBy = append(q.orderBy, []string{string(col), dir.String()})
	return q
}

// Fields injects fields as raw strings into the field clause of the query
// 	sql, e := query.Select(&testassets.Job{}).
// 		Fields(
// 			NewField(FieldTypeBasic, "JobID"),
// 			NewField(FieldTypeBasic, "Name", "Foo"),
// 		)
func (q *Q) Fields(fields ...*Field) *Q {
	q.fields = fields
	return q
}

func (q *Q) Raw(query string) *Q {
	q.raw = query
	return q
}

// Field includes a specific field in the columns to be returned by a result set
func (q *Q) Field(name Column) *Q {

	if _, ok := q.columnTypes[name]; !ok {
		q.errorInvalidColumn(QUERY_ERROR_INVALID_COLUMN, "SELECT...Field", string(name))
	}

	q.fields = append(q.fields, NewField(FieldTypeBasic, name))

	return q
}

// FieldAs includes a specific field in the columns to be returned by a set aliased by `as`
func (q *Q) FieldAs(name Column, as string) *Q {

	if _, ok := q.columnTypes[name]; !ok {
		q.errorInvalidColumn(QUERY_ERROR_INVALID_COLUMN, "SELECT...Field...as", string(name))
	}

	q.fields = append(q.fields, NewField(FieldTypeBasic, name, as))

	return q
}

// FieldRaw allows for an arbitrary string (e.g. "NOW()") to be included in the select columns
func (q *Q) FieldRaw(fieldStr, as string) *Q {
	q.fields = append(q.fields, NewRawField(fieldStr+" AS "+"`"+as+"`"))

	return q
}

// Count creates a count statement
// 	q.Count(query.Column("Foo"), "FooCounted")
// 	COALESCE(COUNT(`t`.`Foo`), 0) AS `FooCounted`
func (q *Q) Count(name Column, as string) *Q {

	if _, ok := q.columnTypes[name]; !ok {
		q.errorInvalidColumn(QUERY_ERROR_INVALID_COLUMN, "SELECT...COUNT()", string(name))
		return q
	}

	q.fields = append(q.fields, NewField(FieldTypeCount, name, as))
	return q

	// return q.FieldRaw("COUNT(`"+q.alias+"`.`"+string(name)+"`)", as)
}

// Sum creates a sum statement
// 	q.Sum(query.Column("Foo"), "FooSummed")
//	COALESCE(SUM(`t`.`Foo`), 0) AS `FooSummed`
func (q *Q) Sum(name Column, as string) *Q {

	if _, ok := q.columnTypes[name]; !ok {
		q.errorInvalidColumn(QUERY_ERROR_INVALID_COLUMN, "SELECT...Sum()", string(name))
		return q
	}

	q.fields = append(q.fields, NewField(FieldTypeSum, name, as))
	return q

	// return q.FieldRaw("COALESCE(SUM(`"+q.alias+"`.`"+string(name)+"`), 0)", as)
}

// Avg creates an Avg statement
// 	q.Avg(query.Column("Foo"), "FooAveraged")
//	COALESCE(AVG(`t`.`Foo`), 0) AS `FooAveraged`
func (q *Q) Avg(name Column, as string) *Q {

	if _, ok := q.columnTypes[name]; !ok {
		q.errorInvalidColumn(QUERY_ERROR_INVALID_COLUMN, "SELECT...Sum()", string(name))
		return q
	}

	q.fields = append(q.fields, NewField(FieldTypeAvg, name, as))
	return q

	// return q.FieldRaw("COALESCE(SUM(`"+q.alias+"`.`"+string(name)+"`), 0)", as)
}

// Min creates a min statement
// 	q.Min(query.Column("Foo"), "MinFoo")
//	COALESCE(MIN(`t`.`Foo`), 0) AS `MinFoo`
func (q *Q) Min(name Column, as string) *Q {

	if _, ok := q.columnTypes[name]; !ok {
		q.errorInvalidColumn(QUERY_ERROR_INVALID_COLUMN, "SELECT...Min()", string(name))
		return q
	}

	q.fields = append(q.fields, NewField(FieldTypeMin, name, as))
	return q

	// return q.FieldRaw("COALESCE(MIN(`"+q.alias+"`.`"+string(name)+"`), 0)", as)
}

// Max creates a max statement
// 	q.Max(query.Column("Foo"), "MaxFoo")
//	COALESCE(MAX(`t`.`Foo`), 0) AS `MaxFoo`
func (q *Q) Max(name Column, as string) *Q {

	if _, ok := q.columnTypes[name]; !ok {
		q.errorInvalidColumn(QUERY_ERROR_INVALID_COLUMN, "SELECT...Max()", string(name))
		return q
	}

	q.fields = append(q.fields, NewField(FieldTypeMax, name, as))
	return q

	// return q.FieldRaw("COALESCE(MAX(`"+q.alias+"`.`"+string(name)+"`), 0)", as)
}

// Where creates or adds to an existing where clause
//
// 	- Simple
// 	q.Where(query.EQ(query.Column("A"), "B"))
// 	WHERE `t`.`B` = `t`.`B`
//
//	- Multiple Arguments
// 	q.Where(query.EQ(query.Column("A"), "B"), query.And(), query.EQ(query.Column("C"), "D"))
//	WHERE `t`.`A` = 'B' AND `t`.`C` = 'D'
//
//	- Daisy Chain
// 	q.Where(query.EQ(query.Column("A"), "B")).Where(query.And()).Where(query.EQ(query.Column("C"), "D"))
//	WHERE `t`.`A` = 'B' AND `t`.`C` = 'D'
//
//	- Separate lines
// 	q.Where(query.EQ(query.Column("A"), "B"))
// 	q.Where(query.And())
// 	q.Where(query.EQ(query.Column("C"), "D"))
//	WHERE `t`.`A` = 'B' AND `t`.`C` = 'D'
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

type QueryErrorType string

const (
	QUERY_ERROR_INVALID_VALUE      QueryErrorType = "Invalid value"
	QUERY_ERROR_INVALID_COLUMN     QueryErrorType = "Invalid Column Name"
	QUERY_ERROR_EMPTY_WHERE_CLAUSE QueryErrorType = "Empty where clause"
)

func (q *Q) String() (string, error) {

	var sb strings.Builder

	switch q.queryType {
	case QueryTypeRaw:
		sb.WriteString(q.raw)
		return sb.String(), nil
	case QueryTypeSelect:

		sb.WriteString("SELECT ")

		if len(q.fields) > 0 {
			for k := range q.fields {
				sb.WriteString(q.FromFieldToString(q.fields[k]))
				if k < len(q.fields)-1 {
					sb.WriteString(", ")
				}
			}
		} else {
			sb.WriteString("`" + q.alias + "`.*")
		}

		sb.WriteString(" FROM")
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
				q.errorInvalidColumn(QUERY_ERROR_INVALID_COLUMN, "UPDATE...SET", string(colName))
				val = "'" + val + "'"
			} else {
				if q.columnTypes[colName] == "%s" {
					val = "'" + EscapeString(val) + "'"
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
				q.errorInvalidColumn(QUERY_ERROR_INVALID_COLUMN, "INSERT...SET", string(colName))
			}

			val := fmt.Sprint(q.sets[colName])

			if q.columnTypes[colName] == "%s" {
				val = "'" + EscapeString(val) + "'"
			}

			cols = append(cols, q.col(string(colName)))
			vals = append(vals, val)
		}

		sb.WriteString(" ( " + strings.Join(cols, ", ") + " ) VALUES ( " + strings.Join(vals, ", ") + " )")
		// sql += " ( " + strings.Join(cols, ", ") + " ) VALUES ( " + strings.Join(vals, ", ") + " )"
	}

	if q.where != nil && q.queryType != QueryTypeInsert {
		whereClause := q.printWhereClause(q.columnTypes, q.where.WhereParts)
		if len(whereClause) > 0 {
			sb.WriteString(" WHERE ")
			sb.WriteString(whereClause)
		}
		// fmt.Println(q.where.WhereParts, q.queryType)
		// q.errorInvalidColumn(QUERY_ERROR_EMPTY_WHERE_CLAUSE, "WHERE", "")
		// q.error(fmt.Sprintf("EMPTY_WHERE_CLAUSE: `%s`", q.model.Table_Name()))

	}

	if q.queryType == QueryTypeSelect && len(q.orderBy) > 0 {
		orderBys := []string{}
		for k := range q.orderBy {

			// Validate the order by column
			if _, ok := q.columnTypes[Column(q.orderBy[k][0])]; !ok {
				q.errorInvalidColumn(QUERY_ERROR_INVALID_COLUMN, "ORDER BY", q.orderBy[k][0])
			}

			orderBys = append(orderBys, q.col(q.orderBy[k][0])+" "+strings.ToUpper(q.orderBy[k][1]))
		}
		sb.WriteString(" ORDER BY " + strings.Join(orderBys, ", "))
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

func (q *Q) printWhereClause(columnTypes map[Column]string, whereParts []*WherePart) string {

	sb := strings.Builder{}

	// prevWasConjunction := false

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
		}

		if w.whereType != WhereTypeExists && w.whereType != WhereTypeNotExists && !isConj && len(w.values) > 0 {

			switch w.whereType {
			case WhereTypeEqualsField:
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

func (q *Q) error(err string) {
	q.errors = append(q.errors, err)
}

func (q *Q) errorInvalidColumn(errType QueryErrorType, queryErrorLocation, comment string) {
	q.error(fmt.Sprintf("%s at %s in model `%s` -- %s", errType, queryErrorLocation, q.model.Table_Name(), comment))
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
// 	wheres := []query.WherePart{
// 		query.EQ(models.ObjectRelationship_Column_IsDeleted, 0),
// 	}
// 	if objectTypeFrom != constants.ObjectTypeUnknown {
// 		wheres = append(wheres, query.And(), query.EQ(models.ObjectRelationship_Column_ObjectTypeFrom, objectTypeFrom))
// 	}
// 	if objectIDFrom > 0 {
// 		wheres = append(wheres, query.And(), query.EQ(models.ObjectRelationship_Column_ObjectIDFrom, objectIDFrom))
// 	}
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

////
// EXPOSED API
////

// EQ is an equals statement between a table column and a value
func EQ(fieldName Column, value interface{}) *WherePart {
	return newWherePart(
		WhereTypeEquals,
		string(fieldName),
		[]interface{}{
			value,
		},
	)
}

// EQF allows for one column to be equal to another
// Example for a subselect
//
// query.Select(&models.UserGroupUser{}).Alias("ugu").FieldRaw("1", "n").Where(
// 	query.EQF("UserID", "`u`.`UserID`"),
// 	query.And(),
// 	query.EQ("UserGroupID", groupID),
// 	query.And(),
// 	query.EQ("IsDeleted", 0),
// ),
func EQF(fieldName1, fieldName2 string) *WherePart {
	return newWherePart(
		WhereTypeEqualsField,
		fieldName1,
		[]interface{}{fieldName2},
	)
}

// NE is a not equals statement between a table column and a value
func NE(fieldName Column, value interface{}) *WherePart {
	return newWherePart(
		WhereTypeNotEquals,
		string(fieldName),
		[]interface{}{value},
	)
}

// LT is a less than statement between a table column and a value
// LT('foo', 1) => WHERE `t`.`foo` < 1
func LT(fieldName Column, value interface{}) *WherePart {
	return newWherePart(
		WhereTypeLessThan,
		string(fieldName),
		[]interface{}{value},
	)
}

// GT is a greater than statement between a table column and a value
func GT(fieldName Column, value interface{}) *WherePart {
	return newWherePart(
		WhereTypeGreaterThan,
		string(fieldName),
		[]interface{}{value},
	)
}

// LTOE is a less than or equals (<=) statement between a table column and a value
//  `t`.`Col` <= value
func LTOE(fieldName Column, value interface{}) *WherePart {
	return newWherePart(
		WhereTypeLessThanOrEqualTo,
		string(fieldName),
		[]interface{}{value},
	)
}

// GTOE is a greater than or equals statement (>=) between a table column and a value
// 	`t`.`Col` >= value
func GTOE(fieldName Column, value interface{}) *WherePart {
	return newWherePart(
		WhereTypeGreaterThanOrEqualTo,
		string(fieldName),
		[]interface{}{value},
	)
}

// Mod is applies modulo operation on column and value testing if it equals remainder
// 	MOD(`t`.`Field`, value) = remainder
func Mod(fieldName Column, value, remainder int64) *WherePart {
	return newWherePart(
		WhereTypeMod,
		string(fieldName),
		[]interface{}{value, remainder},
	)
}

// Modf MOD(value, `t`.`Field`) = remainder
// Example: query.Mod("foo", 2, 1) -> `t`.`Foo` % 2 = 1
func Modf(value int64, fieldName Column, remainder int64) *WherePart {
	return newWherePart(
		WhereTypeModF,
		string(fieldName),
		[]interface{}{value, remainder},
	)
}

// BitAnd `t`.`Field` & a = b
// Example: query.BitAnd("foo", 2, 1) -> `t`.`Foo` & 2 = 1
func BitAnd(fieldName Column, a, b int64) *WherePart {
	return newWherePart(
		WhereTypeBitAnd,
		string(fieldName),
		[]interface{}{a, b},
	)
}

// IN is an IN clause
// Example: query.IN("col1", "foo", "bar", "baz")
func IN(fieldName Column, values ...interface{}) *WherePart {
	return newWherePart(
		WhereTypeIN,
		string(fieldName),
		values,
	)
}

// NOTIN is an NOT IN clause
// Example: query.NOTIN("col1", "foo", "bar", "baz")
// Example: queyr.NOTIN("col2", 1, 2, 3)
func NOTIN(fieldName Column, values ...interface{}) *WherePart {
	return newWherePart(
		WhereTypeNotIN,
		string(fieldName),
		values,
	)
}

// INString is a helper function for converting a slice of string arguments into
// a slice of interface arguments, passed into an IN clause and returned
func INString(fieldName Column, values []string) *WherePart {
	interfaces := make([]interface{}, len(values))

	for k := range values {
		interfaces[k] = values[k]
	}

	return IN(fieldName, interfaces...)
}

// Rawf is a raw SQL statement
// Example: query.Rawf("`t`.`LastRunDate` + 60000 < %d", seconds)),
func Rawf(str string, args ...interface{}) *WherePart {
	return newWherePart(
		WhereTypeRaw,
		"",
		[]interface{}{fmt.Sprintf(str, args...)},
	)
}

// INInt64 is a helper function for converting a slice of string arguments into
// a slice of interface arguments, passed into an IN clause and returned
func INInt64(fieldName Column, values []int64) *WherePart {
	interfaces := make([]interface{}, len(values))

	for k := range values {
		interfaces[k] = values[k]
	}

	return IN(fieldName, interfaces...)
}

// INInt is a helper function for converting a slice of string arguments into
// a slice of interface arguments, passed into an IN clause and returned
func INInt(fieldName Column, values []int) *WherePart {
	interfaces := make([]interface{}, len(values))

	for k := range values {
		interfaces[k] = values[k]
	}

	return IN(fieldName, interfaces...)
}

// Between is a BETWEEN statement
// Example: Between("")
func Between(fieldName Column, from, to interface{}) *WherePart {
	return newWherePart(
		WhereTypeBetween,
		string(fieldName),
		[]interface{}{from, to},
	)
}

func Like(fieldName Column, value string) *WherePart {
	return newWherePart(
		WhereTypeLike,
		string(fieldName),
		[]interface{}{value},
	)
}

func NotLike(fieldName Column, value string) *WherePart {
	return newWherePart(
		WhereTypeNotLike,
		string(fieldName),
		[]interface{}{value},
	)
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
//      SELECT * FROM `Foo` WHERE 1=1 AND FooID = 123;
func WhereAll() *WherePart {
	return newWherePart(
		WhereTypeAll,
		"",
		[]interface{}{},
	)
}

// Exists is a where clause for the SQL EXISTS statement
func Exists(clause *Q) *WherePart {
	clauseString, e := clause.String()

	w := newWherePart(
		WhereTypeExists,
		"",
		[]interface{}{clauseString},
	)
	if e != nil {
		w.e = e
	}
	return w
}

// Exists is a where clause for the SQL EXISTS statement
func NotExists(clause *Q) *WherePart {
	clauseString, e := clause.String()

	w := newWherePart(
		WhereTypeNotExists,
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

func Raw(model IModel, query string) *Q {
	q := Query(model)
	q.queryType = QueryTypeRaw
	q.raw = query
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

func EscapeString(value string) string {
	var sb strings.Builder
	for i := 0; i < len(value); i++ {
		c := value[i]
		switch c {
		case '\\', 0, '\n', '\r', '\'', '"':
			sb.WriteByte('\\')
			sb.WriteByte(c)
		case '\032':
			sb.WriteByte('\\')
			sb.WriteByte('Z')
		default:
			sb.WriteByte(c)
		}
	}
	return sb.String()
}
