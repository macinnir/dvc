package dal

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

// Query defines a query to be made against the database
type Query struct {
	sql   string
	Table *Table

	// Filters are the phrases used in the where clause
	Filters []ValueField

	Params       []string
	Joins        []Join
	resultLimit  int
	resultOffset int
	GroupBy      string
	SetFields    []string
	QueryType    string
	ValueFields  []ValueField
	Dal          *Dal
	Values       []interface{}
	OrderBy      string
	OrderDir     string
	// SelectJoinFields is a collection of fields included in the list of fields selected
	// Format: `joinedTablePrefix`.`fieldName`
	SelectJoinFields []string
}

// ToSQL returns the generated sql string
func (q *Query) ToSQL() string {
	return q.buildSQL()
}

// Query runs the query
func (q *Query) Query() (*sql.Rows, error) {
	query := q.buildSQL()
	return q.Dal.Connection.Query(query, q.Values...)
}

// Exec executes a sql statement
func (q *Query) Exec() (result sql.Result, e error) {
	query := q.buildSQL()
	var stmt *sql.Stmt
	stmt, e = q.Dal.Connection.Prepare(query)
	if e != nil {
		return
	}

	result, e = stmt.Exec(q.Values...)
	return
}

// And adds an and conjuction in the where clause
func (q *Query) And() IQuery {

	q.Where("#and#", "")
	return q
}

// Or adds an or conjuction in the where clause
func (q *Query) Or() IQuery {
	q.Where("#or#", "")
	return q
}

// Where adds a filter on to the where clause
func (q *Query) Where(name string, value interface{}) IQuery {

	q.Filters = append(q.Filters, ValueField{
		Name:  name,
		Value: value,
	})

	return q
}

// Set sets a value for an insert or update statement
func (q *Query) Set(fieldName string, value interface{}) IQuery {

	q.ValueFields = append(q.ValueFields, ValueField{
		Name:  fieldName,
		Value: value,
	})
	return q
}

// Join adds a join clause on to a select statement
func (q *Query) Join(tableName string) IQuery {

	var joinTable *Table
	var ok bool

	if joinTable, ok = q.Dal.Schema.Tables[tableName]; !ok {
		panic("Invaild table name used in join clause")
	}

	q.Joins = append(q.Joins, Join{
		Table:  joinTable,
		Fields: []JoinField{},
	})

	return q
}

// OnValue adds an on clause to the most recent join
func (q *Query) OnValue(fieldName string, value interface{}) IQuery {

	numJoins := len(q.Joins)
	// Get the most recent joins
	if numJoins == 0 {
		panic("Cannot call `On` method before any Joins have been created")
	}

	latestJoinIdx := numJoins - 1

	joinField := JoinField{}
	joinField.FieldName = fieldName
	joinField.Value = value

	q.Joins[latestJoinIdx].Fields = append(q.Joins[latestJoinIdx].Fields, joinField)

	return q
}

// OnField adds an on clause to the most recent join that joins with the field of another table
func (q *Query) OnField(fieldName string, joinTable string, joinField string) IQuery {

	numJoins := len(q.Joins)
	// Get the most recent joins
	if numJoins == 0 {
		panic("Cannot call `On` method before any Joins have been created")
	}

	latestJoinIdx := numJoins - 1

	join := JoinField{}
	join.FieldName = fieldName
	join.JoinTable = joinTable
	join.JoinField = joinField

	q.Joins[latestJoinIdx].Fields = append(q.Joins[latestJoinIdx].Fields, join)

	return q
}

// SelectJoinField uses a field from a joined table in the select list
func (q *Query) SelectJoinField(joinTable string, joinField string, as string) IQuery {
	q.SelectJoinFields = append(q.SelectJoinFields, "`"+q.Dal.Schema.Tables[joinTable].Alias+"`.`"+joinField+"` as `"+as+"`")
	return q
}

// Limit adds a limit clause to the query
func (q *Query) Limit(limit int) IQuery {
	q.resultLimit = limit
	return q
}

// Offset adds an offset clause to the query
func (q *Query) Offset(offset int) IQuery {
	q.resultOffset = offset
	return q
}

// Order adds an orderBy clause to the query
func (q *Query) Order(field string, direction string) IQuery {
	q.OrderBy = field
	q.OrderDir = direction
	return q
}

// GetValues returns the values used for the current query
func (q *Query) GetValues() []interface{} {
	return q.Values
}

func (q *Query) buildInsert() string {

	if len(q.ValueFields) > 0 {
		fields := []string{}
		values := []string{}
		for _, updateField := range q.ValueFields {

			fields = append(fields, "`"+updateField.Name+"`")
			// value := parseValue(updateField.Value)
			q.Values = append(q.Values, updateField.Value)
			values = append(values, "?")
		}
		return fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", q.Table.Name, strings.Join(fields, ","), strings.Join(values, ","))
	}

	panic("No fields to update")
}

func (q *Query) buildUpdate() string {

	setFieldStrings := []string{}

	if len(q.ValueFields) > 0 {

		for _, updateField := range q.ValueFields {

			q.Values = append(q.Values, updateField.Value)

			setFieldStrings = append(setFieldStrings, fmt.Sprintf("`%s`.`%s` = ?", q.Table.Alias, updateField.Name))
		}
	}

	where := q.buildWhere()

	return fmt.Sprintf("UPDATE `%s` `%s` SET %s WHERE %s", q.Table.Name, q.Table.Alias, strings.Join(setFieldStrings, ", "), where)
}

func (q *Query) buildDelete() string {

	where := q.buildWhere()

	return fmt.Sprintf("DELETE FROM `%s` USING `%s` AS `%s` WHERE %s", q.Table.Alias, q.Table.Name, q.Table.Alias, where)
}

func (q *Query) buildCount() string {
	query := "SELECT COUNT(DISTINCT " + fmt.Sprintf("`%s`.`%s`", q.Table.Alias, q.Table.Fields["id"].Name) + ") FROM `" + q.Table.Name + "` `" + q.Table.Alias + "`"

	where := q.buildWhere()

	if len(where) > 0 {
		query = query + " WHERE " + where
	}

	return query
}

func (q *Query) buildSelect() string {

	query := "SELECT "
	colStrings := []string{}
	for _, fk := range q.Table.FieldKeys {
		f := q.Table.Fields[fk]
		colStrings = append(colStrings, fmt.Sprintf("`%s`.`%s`", q.Table.Alias, f.Name))
	}

	if len(q.SelectJoinFields) > 0 {
		for _, jf := range q.SelectJoinFields {
			colStrings = append(colStrings, fmt.Sprintf("%s", jf))
		}
	}

	query = query + strings.Join(colStrings, ", ") + " FROM `" + q.Table.Name + "` `" + q.Table.Alias + "`"

	if len(q.Joins) > 0 {

		joinStrings := ""

		for _, j := range q.Joins {

			onFields := []string{}

			for _, f := range j.Fields {

				onField, eq := q.parseFilterName(f.FieldName)

				if len(f.JoinTable) > 0 && len(f.JoinField) > 0 {
					onField = fmt.Sprintf("`%s`.`%s` %s `%s`.`%s`", j.Table.Alias, onField, eq, q.Dal.Schema.Tables[f.JoinTable].Alias, f.JoinField)
				} else {
					onField = fmt.Sprintf("`%s`.`%s` %s ?", j.Table.Alias, onField, eq)
					q.Values = append(q.Values, f.Value)
				}

				onFields = append(onFields, onField)
			}

			joinString := fmt.Sprintf("JOIN `%s` `%s` ON %s", j.Table.Name, j.Table.Alias, strings.Join(onFields, " AND "))
			joinStrings = joinStrings + " " + joinString
		}

		query = query + joinStrings
	}

	where := q.buildWhere()

	if len(where) > 0 {
		query = query + " WHERE " + where
	}

	if len(q.GroupBy) > 0 {
		query = query + " GROUP BY " + q.GroupBy
	}

	if len(q.OrderBy) > 0 {

		query = query + " ORDER BY " + q.OrderBy

		if len(q.OrderDir) > 0 {
			query = query + " " + q.OrderDir
		} else {
			query = query + " ASC"
		}
	}

	if q.resultLimit > 0 {
		query = query + " LIMIT " + strconv.Itoa(q.resultLimit)
	}

	if q.resultOffset > 0 {
		query = query + " OFFSET " + strconv.Itoa(q.resultOffset)
	}

	return query
}

func (q *Query) parseFilterName(filterName string) (realFilterName string, eq string) {

	eq = "="
	realFilterName = filterName

	filterStartChar := filterName[0:1]

	if filterStartChar == "!" {
		eq = "!="
		realFilterName = filterName[1:]
		return
	}

	if filterStartChar == ">" || filterStartChar == "<" {
		eq = filterStartChar
		realFilterName = filterName[1:]
		return
	}

	if filterName[0:2] == ">=" || filterName[0:2] == "<=" {
		eq = filterName[0:2]
		realFilterName = filterName[2:]
	}

	return
}

func (q *Query) buildWhere() string {

	var eq string
	where := ""
	whereClauses := []string{}
	filterLen := len(q.Filters)
	filterIdx := 0

	if len(q.Filters) > 0 {

		for _, filter := range q.Filters {

			if filter.Name == "#and#" || filter.Name == "#or#" {
				whereClauses = append(whereClauses, strings.ToUpper(filter.Name[1:len(filter.Name)-1]))
				filterIdx = filterIdx + 1
				continue
			}

			filter.Name, eq = q.parseFilterName(filter.Name)

			// value := parseValue(filter.Value)
			q.Values = append(q.Values, filter.Value)
			whereClauses = append(whereClauses, fmt.Sprintf("`%s`.`%s` %s ?", q.Table.Alias, filter.Name, eq))
			filterIdx = filterIdx + 1
			if filterIdx < filterLen {
				if q.Filters[filterIdx].Name[0:1] != "#" {
					whereClauses = append(whereClauses, "AND")
				}
			}
		}
	}

	if len(whereClauses) > 0 {
		where = strings.Join(whereClauses, " ")
	}

	return where
}

//SQL returns the raw sql string
func (q *Query) buildSQL() string {

	if len(q.sql) > 0 {
		return q.sql
	}

	switch q.QueryType {
	case "select":
		q.sql = q.buildSelect()
	case "update":
		q.sql = q.buildUpdate()
	case "delete":
		q.sql = q.buildDelete()
	case "insert":
		q.sql = q.buildInsert()
	case "count":
		q.sql = q.buildCount()
	default:
		panic("Unknown query type")
	}

	// fmt.Printf("SQL: %s\n", sql)

	return q.sql
}

// @TODO Remove
func parseValue(field interface{}) string {

	var value string

	switch field.(type) {
	case int:
		value = strconv.Itoa(field.(int))
	case int64:
		value = strconv.FormatInt(field.(int64), 10)
	case float64:
		value = strconv.FormatFloat(field.(float64), 'E', -1, 64)
	case string:
		value = "'" + field.(string) + "'"
	default:
		panic(fmt.Sprintf("Unknown type for value: %v", field))
	}

	return value
}
