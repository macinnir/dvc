package sqlite

import (
	"bytes"
	"database/sql"
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/macinnir/dvc/lib"
	// sqlite driver
	_ "github.com/mattn/go-sqlite3"
)

// Sqlite implementation of IConnector

// Sqlite contains functionality for interacting with a server
type Sqlite struct {
	Config *lib.Config
}

// Connect connects to a server and returns a new server object
func (ss *Sqlite) Connect() (server *lib.Server, e error) {
	server = &lib.Server{}
	// var connectionString = username + ":" + password + "@tcp(" + host + ")/?charset=utf8"

	server.Connection, e = sql.Open("sqlite3", "./"+ss.Config.Connection.DatabaseName+".db")
	return
}

// FetchDatabases fetches a set of database names from the target server
// populating the Databases property with a map of Database objects
func (ss *Sqlite) FetchDatabases(server *lib.Server) (databases map[string]*lib.Database, e error) {

	// databases in sqlite
	// .databases
	// main: /Users/robertmacinnis/src/github.com/macinnir/dvc/dbTest.db

	// var rows *sql.Rows
	databases = map[string]*lib.Database{}

	databases[ss.Config.Connection.DatabaseName] = &lib.Database{Name: ss.Config.Connection.DatabaseName, Host: ss.Config.Connection.Host}

	return

	// if rows, e = server.Connection.Query("SHOW DATABASES"); e != nil {
	// 	return
	// }

	// if rows != nil {
	// 	defer rows.Close()
	// }

	// for rows.Next() {
	// 	databaseName := ""
	// 	rows.Scan(&databaseName)
	// 	databases[databaseName] = &lib.Database{Name: databaseName, Host: server.Host}
	// }

	// return
}

// UseDatabase switches the connection context to the passed in database
func (ss *Sqlite) UseDatabase(server *lib.Server, databaseName string) (e error) {

	if server.CurrentDatabase == databaseName {
		return
	}

	server.CurrentDatabase = databaseName

	// _, e = server.Connection.Exec(fmt.Sprintf("USE %s", databaseName))
	// if e == nil {
	// 	server.CurrentDatabase = databaseName
	// }
	return
}

// FetchDatabaseTables fetches the complete set of tables from this database
func (ss *Sqlite) FetchDatabaseTables(server *lib.Server, databaseName string) (tables map[string]*lib.Table, e error) {

	tables = map[string]*lib.Table{}

	// var rows *sql.Rows
	// query := "select `TABLE_NAME`, `ENGINE`, `VERSION`, `ROW_FORMAT`, `TABLE_ROWS`, `DATA_LENGTH`, `TABLE_COLLATION`, `AUTO_INCREMENT` FROM information_schema.tables WHERE TABLE_SCHEMA = '" + databaseName + "'"
	// query := ".tables"

	// cmd := exec.Command("sqlite3", "./"+databaseName+".db", ".tables")
	cmd := exec.Command("sqlite3", "./"+databaseName+".db", ".tables")
	var out bytes.Buffer
	cmd.Stdout = &out

	if e = cmd.Run(); e != nil {
		panic(e)
	}

	fields := strings.Fields(out.String())

	if len(fields) > 0 {
		for _, field := range fields {
			table := &lib.Table{Name: field}
			table.Columns, e = ss.FetchTableColumns(server, databaseName, table.Name)
			tables[table.Name] = table
		}
	}

	// fmt.Printf("Output: %s", out.String())

	// var stdOutPipe io.ReadCloser

	// stdOutPipe, e = cmd.StdoutPipe()

	// var stdErr io.ReadCloser

	// if stdErr, e = cmd.StderrPipe(); e != nil {
	// 	panic(e)
	// }

	// if e = cmd.Start(); e != nil {
	// 	panic(e)
	// }

	// errout, _ := ioutil.ReadAll(stdErr)
	// if e = cmd.Wait(); e != nil {
	// 	fmt.Println(errout)
	// 	panic(e)
	// }

	// b, _ := ioutil.ReadAll(stdOutPipe)

	// fmt.Println("Running...")
	// fmt.Println(string(b))

	// 	fmt.Printf("Table: %s\n", table.Name)

	// 	// table.Columns, e = ss.FetchTableColumns(server, databaseName, table.Name)

	// 	if e != nil {
	// 		log.Fatalf("ERROR: %s", e.Error())
	// 		return
	// 	}

	// 	tables[table.Name] = table
	// }

	return
}

type IndexInfo struct {
	ColumnName string
	IndexRank  int64
}

// https://www.sqlite.org/pragma.html#pragma_index_info
func fetchIndexInfo(databaseName string, indexName string) IndexInfo {
	cmd := exec.Command("sqlite3", "./"+databaseName+".db", "PRAGMA index_info("+indexName+")")
	var out bytes.Buffer
	cmd.Stdout = &out
	if e := cmd.Run(); e != nil {
		panic(e)
	}

	rows := strings.Split(out.String(), "\n")

	cols := strings.Split(rows[0], "|")

	i := IndexInfo{
		ColumnName: cols[2],
	}

	i.IndexRank, _ = strconv.ParseInt(cols[0], 0, 64)

	return i
}

type TableIndex struct {
	Name       string
	Unique     bool
	Partial    bool
	ColumnName string
	PrimaryKey bool
}

// Columns:
// 	 1. Ordinal
// 	 2. Name of the index
// 	 3. "1" if UNIQUE, "0" if not
//   4. "c" = created by "CREATE INDEX", "u" if the index was created by "UNIQUE" constraint, "pk" if PRIMARY KEY constraint
//   5. "1" if a partial index, "0" if not
func (ss *Sqlite) FetchTableIndices(databaseName string, tableName string) []TableIndex {

	cmd := exec.Command("sqlite3", "./"+databaseName+".db", "PRAGMA index_list("+tableName+")")
	var out bytes.Buffer
	cmd.Stdout = &out
	if e := cmd.Run(); e != nil {
		panic(e)
	}

	rows := strings.Split(out.String(), "\n")
	columns := []TableIndex{}

	for _, c := range rows {

		if len(c) == 0 {
			continue
		}

		parts := strings.Split(c, "|")

		column := TableIndex{
			Name:       parts[1],
			Unique:     false,
			Partial:    false,
			PrimaryKey: false,
		}

		if parts[2] == "1" {
			column.Unique = true
		}

		if parts[4] == "1" {
			column.Partial = true
		}

		if parts[3] == "pk" {
			column.PrimaryKey = true
		}

		info := fetchIndexInfo(databaseName, column.Name)
		column.ColumnName = info.ColumnName

		columns = append(columns, column)
	}

	return columns
}

// FetchTableColumns lists all of the columns in a table
func (ss *Sqlite) FetchTableColumns(server *lib.Server, databaseName string, tableName string) (columns map[string]*lib.Column, e error) {

	// sqlite3 dbTest.db "PRAGMA table_info(comments)"
	cmd := exec.Command("sqlite3", "./"+databaseName+".db", "PRAGMA table_info("+tableName+")")
	var out bytes.Buffer
	cmd.Stdout = &out
	if e = cmd.Run(); e != nil {
		panic(e)
	}
	columnsData := strings.Split(out.String(), "\n")

	columns = map[string]*lib.Column{}

	indices := ss.FetchTableIndices(databaseName, tableName)

	findColumnIndex := func(columnName string) TableIndex {

		t := TableIndex{}

		for _, i := range indices {
			if i.ColumnName == columnName {
				t = i
			}
		}

		return t
	}

	for position, c := range columnsData {
		if len(c) == 0 {
			continue
		}

		key := ""

		parts := strings.Split(c, "|")

		if parts[5] == "1" {
			key = "PRI"
		}

		column := lib.Column{
			Name:       parts[1],
			Position:   position + 1,
			Type:       parts[2],
			DataType:   parts[2],
			IsNullable: parts[3] == "0",
			ColumnKey:  key,
			Default:    parts[4],
		}

		idx := findColumnIndex(column.Name)
		if len(idx.ColumnName) > 0 {
			if idx.Unique == true {
				column.ColumnKey = "UNI"
			} else {
				column.ColumnKey = "MUL"
			}
		}

		columns[parts[1]] = &column
	}

	return
}

// CreateChangeSQL generates sql statements based off of comparing two database objects
// localSchema is authority, remoteSchema will be upgraded to match localSchema
func (ss *Sqlite) CreateChangeSQL(localSchema *lib.Database, remoteSchema *lib.Database) (sql string, e error) {

	query := ""

	// What tables are in local that aren't in remote?
	for tableName, table := range localSchema.Tables {

		// Table does not exist on remote schema
		if _, ok := remoteSchema.Tables[tableName]; !ok {

			// fmt.Printf("Local table %s is not in remote\n", table.Name)
			query, e = ss.createTable(table)
			// fmt.Printf("Running Query: %s\n", query)
			sql += query + "\n"
		} else {
			remoteTable := remoteSchema.Tables[tableName]
			query, e = ss.createTableChangeSQL(table, remoteTable)
			if len(query) > 0 {
				sql += query + "\n"
			}
		}
	}

	// What tables are in remote that aren't in local?
	for _, table := range remoteSchema.Tables {

		// Table does not exist on local schema
		if _, ok := localSchema.Tables[table.Name]; !ok {
			query, e = ss.dropTable(table)
			sql += query + "\n"
		}
	}

	return
}

// createTableChangeSQL returns a set of statements that alter a table's structure if and only if there is a difference between
// the local and remote tables
// If no change is found, an empty string is returned.
func (ss *Sqlite) createTableChangeSQL(localTable *lib.Table, remoteTable *lib.Table) (sql string, e error) {

	var query string

	for _, column := range localTable.Columns {

		// Column does not exist remotely
		if _, ok := remoteTable.Columns[column.Name]; !ok {
			query, e = ss.alterTableCreateColumn(localTable, column)
			if e != nil {
				return
			}

			if len(query) > 0 {
				sql += query + "\n"
			}

		} else {

			remoteColumn := remoteTable.Columns[column.Name]

			query, e = ss.changeColumn(localTable, column, remoteColumn)

			if e != nil {
				return
			}

			if len(query) > 0 {
				sql += query + "\n"
			}
		}
	}

	for _, column := range remoteTable.Columns {

		// Column does not exist locally
		if _, ok := localTable.Columns[column.Name]; !ok {
			query, e = ss.alterTableDropColumn(localTable, column)
			if e != nil {
				return
			}

			sql += query + "\n"
		}
	}

	return
}

// createTable returns a create table sql statement
func (ss *Sqlite) createTable(table *lib.Table) (sql string, e error) {

	// colLen := len(table.Columns)
	idx := 1

	cols := []string{}

	// Unique Keys
	uniqueKeyColumns := []*lib.Column{}

	// Regular Keys (allows for multiple entries)
	multiKeyColumns := []*lib.Column{}

	sortedColumns := make(lib.SortedColumns, 0, len(table.Columns))

	for _, column := range table.Columns {
		sortedColumns = append(sortedColumns, column)
	}

	sort.Sort(sortedColumns)

	for _, column := range sortedColumns {

		colQuery := ""
		colQuery, e = ss.createColumn(column)
		col := colQuery

		idx++

		switch column.ColumnKey {
		case "PRI":
			col += " PRIMARY KEY AUTOINCREMENT"
		case "UNI":
			uniqueKeyColumns = append(uniqueKeyColumns, column)
		case "MUL":
			multiKeyColumns = append(multiKeyColumns, column)
		}
		cols = append(cols, col)
	}

	sql = fmt.Sprintf("CREATE TABLE %s (\n\t%s\n);", table.Name, strings.Join(cols, ",\n\t"))

	if len(uniqueKeyColumns) > 0 {
		sql += "\n"
		for _, uniqueKeyColumn := range uniqueKeyColumns {
			t, _ := ss.addUniqueIndex(table, uniqueKeyColumn)
			sql += t + "\n"
		}
	}

	if len(multiKeyColumns) > 0 {
		sql += "\n"
		for _, multiKeyColumn := range multiKeyColumns {
			t, _ := ss.addIndex(table, multiKeyColumn)
			sql += t + "\n"
		}
	}

	return
}

// dropTable returns a drop table sql statement
func (ss *Sqlite) dropTable(table *lib.Table) (sql string, e error) {
	sql = fmt.Sprintf("DROP TABLE \"%s\";", table.Name)
	return
}

func (ss *Sqlite) translateColumnType(column *lib.Column) (sqlType string) {

	switch strings.ToLower(column.DataType) {
	case "bigint":
		sqlType = "INTEGER"
	case "smallint":
		sqlType = "INTEGER"
	case "mediumint":
		sqlType = "INTEGER"
	case "varchar":
		sqlType = "TEXT"
	case "char":
		sqlType = "TEXT"
	case "text":
		sqlType = "TEXT"
	case "datetime":
		sqlType = "TEXT"
	case "decimal":
		sqlType = "NUMERIC"
	case "float":
		sqlType = "NUMERIC"
	case "enum":
		sqlType = "TEXT"
	case "int":
		sqlType = "INTEGER"
	case "tinyint":
		sqlType = "INTEGER"
	default:
		fmt.Printf("Nerp: %s\n", column.DataType)
	}
	return
}

// createColumn returns a table column sql segment
func (ss *Sqlite) createColumn(column *lib.Column) (sql string, e error) {

	dt := ss.translateColumnType(column)

	sql = fmt.Sprintf("%s %s", column.Name, dt)
	if !column.IsNullable {
		sql += " NOT"
	}
	sql += " NULL"

	if !column.IsNullable && (dt == "TEXT") {
		sql += fmt.Sprintf(" DEFAULT '%s'", column.Default)
	} else if len(column.Default) > 0 {
		sql += fmt.Sprintf(" DEFAULT %s", column.Default)
	}

	// if len(column.Extra) > 0 {
	// 	sql += " " + column.Extra
	// }

	return

}

// alterTableDropColumn returns an alter table sql statement that drops a column
func (ss *Sqlite) alterTableDropColumn(table *lib.Table, column *lib.Column) (sql string, e error) {
	sql = fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN `%s`;", table.Name, column.Name)
	return
}

// changeColumn returns an alter table sql statement that adds or removes an index from a column
// if and only if the one (e.g. local) has a column and the other (e.g. remote) does not
// Truth table
// 		Remote 	| 	Local 	| 	Result
// ---------------------------------------------------------
// 1. 	MUL		| 	none 	| 	Drop index
// 2. 	UNI		| 	none 	| 	Drop unique index
// 3. 	none 	| 	MUL 	|  	Create index
// 4. 	none 	| 	UNI 	| 	Create unique index
// 5. 	MUL		| 	UNI 	| 	Drop index; Create unique index
// 6. 	UNI 	| 	MUL 	| 	Drop unique index; Create index
// 7. 	none	| 	none	| 	Do nothing
// 8. 	MUL		| 	MUL		| 	Do nothing
// 9. 	UNI		|   UNI		| 	Do nothing
func (ss *Sqlite) changeColumn(table *lib.Table, localColumn *lib.Column, remoteColumn *lib.Column) (sql string, e error) {

	t := ""
	query := ""

	// 7,8,9
	if localColumn.ColumnKey == remoteColumn.ColumnKey {
		return
	}

	// <7
	if localColumn.ColumnKey != remoteColumn.ColumnKey {

		// 1,2: There is no indexing on the local schema
		if localColumn.ColumnKey == "" {
			switch remoteColumn.ColumnKey {
			// 1
			case "MUL":
				t, _ = ss.dropIndex(table, localColumn)
				query += t + "\n"
			// 2
			case "UNI":
				t, _ = ss.dropUniqueIndex(table, localColumn)
				query += t + "\n"
			}
		}

		// 3, 4: There is no indexing on the remote schema
		if remoteColumn.ColumnKey == "" {
			switch localColumn.ColumnKey {
			// 3
			case "MUL":
				t, _ = ss.addIndex(table, localColumn)
				query += t + "\n"
			// 4
			case "UNI":
				t, _ = ss.addUniqueIndex(table, localColumn)
				query += t + "\n"
			}
		}

		// 5
		if remoteColumn.ColumnKey == "MUL" && localColumn.ColumnKey == "UNI" {
			t, _ = ss.dropIndex(table, localColumn)
			query += t + "\n"
			t, _ = ss.addUniqueIndex(table, localColumn)
			query += t + "\n"
		}

		// 6
		if remoteColumn.ColumnKey == "UNI" && localColumn.ColumnKey == "MUL" {
			t, _ = ss.dropUniqueIndex(table, localColumn)
			query += t + "\n"
			t, _ = ss.addIndex(table, localColumn)
			query += t + "\n"
		}
	}

	sql = query
	return

}

// alterTableCreateColumn returns an alter table sql statement that adds a column
func (ss *Sqlite) alterTableCreateColumn(table *lib.Table, column *lib.Column) (sql string, e error) {

	query := ""

	query, e = ss.createColumn(column)
	sql = fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s;", table.Name, query)

	return
}

// addIndex returns an alter table sql statement that adds an index to a table
func (ss *Sqlite) addIndex(table *lib.Table, column *lib.Column) (sql string, e error) {
	sql = fmt.Sprintf("CREATE INDEX i_%s_%s ON %s (%s);", table.Name, column.Name, table.Name, column.Name)
	return
}

// addUniqueIndex returns an alter table sql statement that adds a unique index to a table
func (ss *Sqlite) addUniqueIndex(table *lib.Table, column *lib.Column) (sql string, e error) {
	sql = fmt.Sprintf("CREATE UNIQUE INDEX ui_%s_%s ON %s (%s);", table.Name, column.Name, table.Name, column.Name)
	return
}

// dropIndex returns an alter table sql statement that drops an index
func (ss *Sqlite) dropIndex(table *lib.Table, column *lib.Column) (sql string, e error) {
	sql = fmt.Sprintf("DROP INDEX i_%s_%s;", table.Name, column.Name)
	return
}

// dropUniqueIndex returns an alter table sql statement that drops a unique index
func (ss *Sqlite) dropUniqueIndex(table *lib.Table, column *lib.Column) (sql string, e error) {
	sql = fmt.Sprintf("DROP INDEX ui_%s_%s;", table.Name, column.Name)
	return
}
