package dal

import "database/sql"

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
