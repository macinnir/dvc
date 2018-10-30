package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/macinnir/dvc/lib"
)

// MySQL implementation of IConnector

// MySQL contains functionality for interacting with a server
type MySQL struct {
	Config *lib.Config
}

// Connect connects to a server and returns a new server object
func (ss *MySQL) Connect() (server *lib.Server, e error) {
	server = &lib.Server{Host: ss.Config.Host}
	var connectionString = ss.Config.Username + ":" + ss.Config.Password + "@tcp(" + ss.Config.Host + ")/?charset=utf8"
	server.Connection, e = sql.Open("mysql", connectionString)
	return
}

// FetchDatabases fetches a set of database names from the target server
// populating the Databases property with a map of Database objects
func (ss *MySQL) FetchDatabases(server *lib.Server) (databases map[string]*lib.Database, e error) {

	var rows *sql.Rows
	databases = map[string]*lib.Database{}

	if rows, e = server.Connection.Query("SHOW DATABASES"); e != nil {
		return
	}

	if rows != nil {
		defer rows.Close()
	}

	for rows.Next() {
		databaseName := ""
		rows.Scan(&databaseName)
		databases[databaseName] = &lib.Database{Name: databaseName, Host: server.Host}
	}

	return
}

// UseDatabase switches the connection context to the passed in database
func (ss *MySQL) UseDatabase(server *lib.Server, databaseName string) (e error) {

	if server.CurrentDatabase == databaseName {
		return
	}

	_, e = server.Connection.Exec(fmt.Sprintf("USE %s", databaseName))
	if e == nil {
		server.CurrentDatabase = databaseName
	}
	return
}

// FetchDatabaseTables fetches the complete set of tables from this database
func (ss *MySQL) FetchDatabaseTables(server *lib.Server, databaseName string) (tables map[string]*lib.Table, e error) {

	var rows *sql.Rows
	query := "select `TABLE_NAME`, `ENGINE`, `VERSION`, `ROW_FORMAT`, `TABLE_ROWS`, `DATA_LENGTH`, `TABLE_COLLATION`, `AUTO_INCREMENT` FROM information_schema.tables WHERE TABLE_SCHEMA = '" + databaseName + "'"
	// fmt.Printf("Query: %s\n", query)
	if rows, e = server.Connection.Query(query); e != nil {
		return
	}

	if rows != nil {
		defer rows.Close()
	}

	tables = map[string]*lib.Table{}

	for rows.Next() {

		table := &lib.Table{}

		rows.Scan(
			&table.Name,
			&table.Engine,
			&table.Version,
			&table.RowFormat,
			&table.Rows,
			&table.DataLength,
			&table.Collation,
			&table.AutoIncrement,
		)

		table.Columns, e = ss.FetchTableColumns(server, databaseName, table.Name)

		if e != nil {
			log.Fatalf("ERROR: %s", e.Error())
			return
		}

		tables[table.Name] = table
	}

	return
}

// FetchTableColumns lists all of the columns in a table
func (ss *MySQL) FetchTableColumns(server *lib.Server, databaseName string, tableName string) (columns map[string]*lib.Column, e error) {

	var rows *sql.Rows

	query := fmt.Sprintf(`
		SELECT 
			COLUMN_NAME, 
			ORDINAL_POSITION, 
			COALESCE(COLUMN_DEFAULT, '') as COLUMN_DEFAULT, 
			CASE IS_NULLABLE 
				WHEN 'YES' THEN 1
				ELSE 0
			END AS IS_NULLABLE,
			DATA_TYPE, 
			COALESCE(CHARACTER_MAXIMUM_LENGTH, 0) as MaxLength, 
			COALESCE(NUMERIC_PRECISION, 0) as NumericPrecision, 
			COALESCE(CHARACTER_SET_NAME, '') AS CharSet, 
			COLUMN_TYPE,
			COLUMN_KEY,
			EXTRA,
			COALESCE(NUMERIC_SCALE, 0) as NumericScale 
		FROM information_schema.COLUMNS 
		WHERE 
			TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s'
	`, databaseName, tableName)

	if rows, e = server.Connection.Query(query); e != nil {
		return
	}

	if rows != nil {
		defer rows.Close()
	}

	columns = map[string]*lib.Column{}

	for rows.Next() {
		column := lib.Column{}
		if e = rows.Scan(
			&column.Name,
			&column.Position,
			&column.Default,
			&column.IsNullable,
			&column.DataType,
			&column.MaxLength,
			&column.Precision,
			&column.CharSet,
			&column.Type,
			&column.ColumnKey,
			&column.Extra,
			&column.NumericScale,
		); e != nil {
			return
		}
		column.IsUnsigned = strings.Contains(strings.ToLower(column.Type), " unsigned")
		columns[column.Name] = &column
	}

	// fmt.Printf("Fetching columns database: %s, table: %s - columns: %d\n", databaseName, tableName, len(columns))

	return
}

// CreateChangeSQL generates sql statements based off of comparing two database objects
// localSchema is authority, remoteSchema will be upgraded to match localSchema
func (ss *MySQL) CreateChangeSQL(localSchema *lib.Database, remoteSchema *lib.Database) (sql string, e error) {

	fmt.Println("Creating MySQL change SQL!!!!")

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
func (ss *MySQL) createTableChangeSQL(localTable *lib.Table, remoteTable *lib.Table) (sql string, e error) {

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
func (ss *MySQL) createTable(table *lib.Table) (sql string, e error) {

	// colLen := len(table.Columns)
	idx := 1

	// Primary Key?
	primaryKey := ""

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
			primaryKey = column.Name
		case "UNI":
			uniqueKeyColumns = append(uniqueKeyColumns, column)
		case "MUL":
			multiKeyColumns = append(multiKeyColumns, column)
		}
		cols = append(cols, col)
	}

	if len(primaryKey) > 0 {
		cols = append(cols, fmt.Sprintf("PRIMARY KEY(`%s`)", primaryKey))
	}

	sql = fmt.Sprintf("CREATE TABLE `%s` (\n\t%s\n) ENGINE = %s;", table.Name, strings.Join(cols, ",\n\t"), table.Engine)

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
func (ss *MySQL) dropTable(table *lib.Table) (sql string, e error) {
	sql = fmt.Sprintf("DROP TABLE `%s`;", table.Name)
	return
}

// isInt
// Integer DataTypes: https://dev.mysql.com/doc/refman/8.0/en/integer-types.html
func (ss *MySQL) isInt(dataType string) bool {
	switch strings.ToLower(dataType) {
	case "tinyint":
		return true
	case "smallint":
		return true
	case "mediumint":
		return true
	case "int":
		return true
	case "bigint":
		return true
	}
	return false
}

func (ss *MySQL) intColLength(dataType string, isUnsigned bool) int {
	switch dataType {
	case "tinyint":
		if isUnsigned {
			return 3
		}
		return 4
	case "smallint":
		if isUnsigned {
			return 5
		}
		return 6
	case "mediumint":
		if isUnsigned {
			return 8
		}
		return 9
	case "int":
		if isUnsigned {
			return 10
		}
		return 11
	case "bigint":
		return 20
	}

	return 0
}

// Fixed Point Types
// https://dev.mysql.com/doc/refman/8.0/en/fixed-point-types.html
func (ss *MySQL) isFixedPointType(dataType string) bool {
	switch strings.ToLower(dataType) {
	case "decimal":
		return true
	case "numeric":
		return true
	}
	return false
}

// Floating Point Types
// https://dev.mysql.com/doc/refman/8.0/en/floating-point-types.html
func (ss *MySQL) isFloatingPointType(dataType string) bool {
	switch strings.ToLower(dataType) {
	case "float":
		return true
	case "double":
		return true
	}

	return false
}

func (ss *MySQL) isString(dataType string) bool {
	switch strings.ToLower(dataType) {
	case "varchar":
		return true
	case "char":
		return true
	}

	return false
}

func (ss *MySQL) hasDefaultString(dataType string) bool {
	switch strings.ToLower(dataType) {
	case "varchar":
		return true
	case "char":
		return true
	case "enum":
		return true
	}
	return false
}

// createColumn returns a table column sql segment
// Data Types
// INT SIGNED 	11 columns
//

func (ss *MySQL) createColumn(column *lib.Column) (sql string, e error) {

	if ss.isInt(column.DataType) {

		sql = fmt.Sprintf("`%s` %s(%d)", column.Name, column.DataType, ss.intColLength(column.DataType, column.IsUnsigned))

		if column.IsUnsigned == true {
			sql += " UNSIGNED "
		} else {
			sql += " SIGNED "
		}

	} else if ss.isFixedPointType(column.DataType) {
		sql = fmt.Sprintf("`%s` %s(%d,%d)", column.Name, column.DataType, column.Precision, column.NumericScale)
		if column.IsUnsigned == true {
			sql += " UNSIGNED "
		} else {
			sql += " SIGNED "
		}
	} else if ss.isFloatingPointType(column.DataType) {
		sql = fmt.Sprintf("`%s` %s(%d,%d)", column.Name, column.DataType, column.Precision, column.NumericScale)
		if column.IsUnsigned == true {
			sql += " UNSIGNED "
		} else {
			sql += " SIGNED "
		}
	} else if ss.isString(column.DataType) {
		sql = fmt.Sprintf("`%s` %s(%d)", column.Name, column.DataType, column.MaxLength)
	} else {
		sql = fmt.Sprintf("`%s` %s", column.Name, column.DataType)
	}

	if !column.IsNullable {
		sql += " NOT"
	}
	sql += " NULL"

	// Add single quotes to string default
	if ss.hasDefaultString(column.DataType) {
		sql += fmt.Sprintf(" DEFAULT '%s'", column.Default)
	} else if len(column.Default) > 0 {
		sql += fmt.Sprintf(" DEFAULT %s", column.Default)
	}

	if len(column.Extra) > 0 {
		sql += " " + column.Extra
	}

	return

}

// alterTableDropColumn returns an alter table sql statement that drops a column
func (ss *MySQL) alterTableDropColumn(table *lib.Table, column *lib.Column) (sql string, e error) {
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
func (ss *MySQL) changeColumn(table *lib.Table, localColumn *lib.Column, remoteColumn *lib.Column) (sql string, e error) {

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
func (ss *MySQL) alterTableCreateColumn(table *lib.Table, column *lib.Column) (sql string, e error) {

	query := ""

	query, e = ss.createColumn(column)
	sql = fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN %s;", table.Name, query)

	return
}

// addIndex returns an alter table sql statement that adds an index to a table
func (ss *MySQL) addIndex(table *lib.Table, column *lib.Column) (sql string, e error) {
	sql = fmt.Sprintf("ALTER TABLE `%s` ADD INDEX `i_%s` (`%s`);", table.Name, column.Name, column.Name)
	return
}

// addUniqueIndex returns an alter table sql statement that adds a unique index to a table
func (ss *MySQL) addUniqueIndex(table *lib.Table, column *lib.Column) (sql string, e error) {
	sql = fmt.Sprintf("ALTER TABLE `%s` ADD UNIQUE INDEX `ui_%s` (`%s`);", table.Name, column.Name, column.Name)
	return
}

// dropIndex returns an alter table sql statement that drops an index
func (ss *MySQL) dropIndex(table *lib.Table, column *lib.Column) (sql string, e error) {
	sql = fmt.Sprintf("ALTER TABLE `%s` DROP INDEX `i_%s`;", table.Name, column.Name)
	return
}

// dropUniqueIndex returns an alter table sql statement that drops a unique index
func (ss *MySQL) dropUniqueIndex(table *lib.Table, column *lib.Column) (sql string, e error) {
	sql = fmt.Sprintf("ALTER TABLE `%s` DROP INDEX `ui_%s`;", table.Name, column.Name)
	return
}

// // FetchDatabase builds and fetches data for a database object
// func FetchDatabase(server *Server, databaseName string) (database *lib.Database, e error) {

// 	database = &Database{
// 		Host: server.Host,
// 		Name: databaseName,
// 	}

// 	database.Tables, e = FetchTables(server, databaseName)

// 	return
// }

// // CreateDatabase creates a new databse
// func (s *Server) CreateDatabase(databaseName string) (database *lib.Database, e error) {
// 	_, e = s.conn.Exec(fmt.Sprintf("CREATE DATABASE `%s`", databaseName))
// 	if e != nil {
// 		return
// 	}
// 	s.Databases[databaseName] = &Database{name: databaseName, host: s.Name}
// 	database = s.Databases[databaseName]
// 	return
// }

// // NewServer creates a new instance of a Server object
// func NewServer(host string, username string, password string) (s Server) {
// 	s = Server{Name: host}
// 	e := s.Connect(username, password)
// 	if e != nil {
// 		log.Fatal(e)
// 	}

// 	return
// }
