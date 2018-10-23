package mysql

import (
	"database/sql"
	"fmt"
	"log"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/macinnir/dvc/lib"
)

// MySQL implementation of IConnector

// MySQL contains functionality for interacting with a server
type MySQL struct {
	Config *lib.Config
}

// ConnectToServer connects to a server and returns a new server object
func (ss *MySQL) ConnectToServer(host string, username string, password string) (server *lib.Server, e error) {
	server = &lib.Server{Host: host}
	var connectionString = username + ":" + password + "@tcp(" + host + ")/?charset=utf8"
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
			EXTRA
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
		); e != nil {
			return
		}
		columns[column.Name] = &column
	}

	// fmt.Printf("Fetching columns database: %s, table: %s - columns: %d\n", databaseName, tableName, len(columns))

	return
}

// // FetchDatabase builds and fetches data for a database object
// func FetchDatabase(server *Server, databaseName string) (database *Database, e error) {

// 	database = &Database{
// 		Host: server.Host,
// 		Name: databaseName,
// 	}

// 	database.Tables, e = FetchTables(server, databaseName)

// 	return
// }

// // CreateDatabase creates a new databse
// func (s *Server) CreateDatabase(databaseName string) (database *Database, e error) {
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
