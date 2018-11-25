package dal

import (
	"database/sql"
	"fmt"
)

// Field represents a field on a table
type Field struct {
	Name         string
	Type         string
	DefaultValue string
}

// NewTable creates a new table object
func NewTable(name string) *Table {
	t := new(Table)
	t.Name = name
	t.Fields = map[string]*Field{}
	return t
}

// Table represents a database table
type Table struct {
	Name      string
	Fields    map[string]*Field
	Alias     string
	FieldKeys []string
}

// AddFields adds a collection of fields
func (t *Table) AddFields(fieldNames []string) {
	for _, fieldName := range fieldNames {
		t.AddField(fieldName)
	}
}

// AddField adds a new field to the database object
func (t *Table) AddField(fieldName string) (e error) {

	var ok bool

	if _, ok = t.Fields[fieldName]; ok {
		e = fmt.Errorf("the field `%s`.`%s` has already been defined", t.Name, fieldName)
		return
	}

	newField := new(Field)
	newField.Name = fieldName
	t.Fields[fieldName] = newField
	t.FieldKeys = append(t.FieldKeys, fieldName)

	return
}

// Field gets a field from the table by its name
func (t *Table) Field(fieldName string) *Field {

	var field *Field
	var ok bool

	if field, ok = t.Fields[fieldName]; !ok {
		panic(fmt.Sprintf("Field `%s`.`%s` not found", t.Name, fieldName))
	}

	return field
}

// ValueField is a name/value pair used to setting data on insert or update queries
type ValueField struct {
	Name  string
	Value interface{}
}

// JoinField represents a part of a join clause
type JoinField struct {
	FieldName string
	Value     interface{}
	JoinTable string
	JoinField string
}

// Join represents a join clause
type Join struct {
	Table  *Table
	Fields []JoinField
}

// IQuery outlines the methods on build a sql query and interacting with the database
type IQuery interface {
	Query() (*sql.Rows, error)
	Exec() (result sql.Result, e error)
	And() IQuery
	Or() IQuery
	Where(name string, value interface{}) IQuery
	Set(fieldName string, value interface{}) IQuery
	Join(tableName string) IQuery
	OnValue(tableName string, value interface{}) IQuery
	OnField(fieldName string, joinTable string, joinField string) IQuery
	Limit(limit int) IQuery
	Offset(offset int) IQuery
	Order(field string, direction string) IQuery
	ToSQL() string
	GetValues() []interface{}
	SelectJoinField(tableName string, fieldName string, as string) IQuery
}

// ISchema represents DAL schema methods
type ISchema interface {
	Select(tableName string) IQuery
	Update(tableName string) IQuery
	Delete(tableName string) IQuery
	Insert(tableName string) IQuery
	AddTable(name string, fields []string) error
	Table(name string) (t *Table)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (result *sql.Rows, e error)
	GetTables() map[string]*Table
}
