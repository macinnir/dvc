package main

import (
	"database/sql"
	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// MySQL integration
type MySQL struct {
	conn   *sql.DB
	config *Config
}

func NewMySQL(config *Config) (mysql *MySQL, e error) {

	mysql = &MySQL{
		config: config,
	}

	return
}

// Connect connects to the database
func (d *MySQL) Connect() (e error) {
	var connectionString = d.config.Username + ":" + d.config.Password + "@tcp(" + d.config.Host + ")/" + d.config.DatabaseName + "?charset=utf8"
	d.conn, e = sql.Open("mysql", connectionString)
	return
}

// CreateRun creates a new dvcRun
func (d *MySQL) CreateRun() (runID int64, e error) {
	var result sql.Result
	result, e = d.conn.Exec("INSERT INTO dvcRun (`dateCreated`) VALUES (UNIX_TIMESTAMP())")
	if e != nil {
		return
	}
	runID, e = result.LastInsertId()
	return
}

// FetchAllTableNames fetches the complete set of tables from this database
// and populates the Tables map with a collection of Table objects
func (d *MySQL) FetchAllTableNames() (tableNames map[string]string, e error) {

	var rows *sql.Rows

	if rows, e = d.conn.Query("SHOW TABLES"); e != nil {
		return
	}

	if rows != nil {
		defer rows.Close()
	}

	tableNames = map[string]string{}
	var tableName string

	for rows.Next() {
		rows.Scan(&tableName)
		tableNames[tableName] = tableName
	}

	return
}

// CreateBaseTablesIfNotExists creates the base tables
func (d *MySQL) CreateBaseTablesIfNotExists() (e error) {
	return
}

// DropBaseTables drops the base tables
func (d *MySQL) DropBaseTables() (e error) {
	return
}

// func FetchFilesBySet(dvcSetID int64, dvcSetName string) (dvcFiles map[string]*DVCFile, sortedFileKeys []string, e error) {
// 	return
// }
// func FetchSets() (e error) {
// 	return
// }
// func CreateSet(name string) (e error) {
// 	return
// }
// func DeleteSet(setName string) (e error) {
// 	return
// }
// func CreateFile(set *DVCSet, name string, fileHash string, content string) (f *DVCFile, e error) {
// 	return
// }
// func DeleteFile(set *DVCSet, file *DVCFile) (e error) {
// 	return
// }
// func FetchTables() (e error) {
// 	return
// }
// func SetExists(changesetName string) (ok bool) {
// 	return
// }
// func FileExists(changesetName string, fileName string) (ok bool) {
// 	return
// }
// func TableExists(tableName string) (ok bool) {
// 	return
// }
// func RunChange(sqlString string) (e error) {
// 	return
// }
// func SetFileAsRun(dvcFileId int64) (e error) {
// 	return
// }
