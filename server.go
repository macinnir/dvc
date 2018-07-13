package main

import (
	"database/sql"
	"fmt"
	"log"
	// "fmt"
	// "log"
)

// serverService contains functionality for interacting with a server
type serverService struct {
	Config *Config
}

// ConnectToServer connects to a server and returns a new server object
func (ss *serverService) ConnectToServer(host string, username string, password string) (server *Server, e error) {
	server = &Server{Host: host}
	var connectionString = username + ":" + password + "@tcp(" + host + ")/?charset=utf8"
	server.Connection, e = sql.Open("mysql", connectionString)
	return
}

// FetchDatabases fetches a set of database names from the target server
// populating the Databases property with a map of Database objects
func (ss *serverService) FetchDatabases(server *Server) (databases map[string]*Database, e error) {

	var rows *sql.Rows
	databases = map[string]*Database{}

	if rows, e = server.Connection.Query("SHOW DATABASES"); e != nil {
		return
	}

	if rows != nil {
		defer rows.Close()
	}

	for rows.Next() {
		databaseName := ""
		rows.Scan(&databaseName)
		databases[databaseName] = &Database{Name: databaseName, Host: server.Host}
	}

	return
}

// UseDatabase switches the connection context to the passed in database
func (ss *serverService) UseDatabase(server *Server, databaseName string) (e error) {

	if server.CurrentDatabase == databaseName {
		return
	}

	_, e = server.Connection.Exec(fmt.Sprintf("USE %s", databaseName))
	if e != nil {
		server.CurrentDatabase = databaseName
	}
	return
}

// FetchDatabaseTables fetches the complete set of tables from this database
func (ss *serverService) FetchDatabaseTables(server *Server, databaseName string) (tables map[string]*Table, e error) {

	var rows *sql.Rows
	query := "select `TABLE_NAME`, `ENGINE`, `VERSION`, `ROW_FORMAT`, `TABLE_ROWS`, `DATA_LENGTH`, `TABLE_COLLATION`, `AUTO_INCREMENT` FROM information_schema.tables WHERE TABLE_SCHEMA = '" + databaseName + "'"
	// fmt.Printf("Query: %s\n", query)
	if rows, e = server.Connection.Query(query); e != nil {
		return
	}

	if rows != nil {
		defer rows.Close()
	}

	tables = map[string]*Table{}

	for rows.Next() {
		table := &Table{}

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

		table.Columns, e = FetchColumns(server, table.Name)

		if e != nil {
			log.Fatalf("ERROR: %s", e.Error())
			return
		}

		tables[table.Name] = table
	}

	return
}

// FetchTableColumns lists all of the columns in a table
func (ss *serverService) FetchTableColumns(server *Server, tableName string) (columns map[string]*Column, e error) {
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
	`, server.CurrentDatabase, tableName)

	if rows, e = server.Connection.Query(query); e != nil {
		return
	}

	if rows != nil {
		defer rows.Close()
	}

	columns = map[string]*Column{}

	for rows.Next() {
		column := Column{}
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

	return
}

// Server represents a server
type Server struct {
	Name            string `json:"name"`
	Host            string `json:"host"`
	Databases       map[string]*Database
	Connection      *sql.DB
	CurrentDatabase string
}

// DatabaseExists checks if the database `databaseName` exists in its list of databases
func (s *Server) DatabaseExists(databaseName string) bool {

	exists := false

	for _, db := range s.Databases {
		if db.Name == databaseName {
			exists = true
			break
		}
	}

	return exists
}

// Database represents a database
type Database struct {
	RunID int64
	Name  string
	Host  string
	// Sets          map[string]*ChangeSet
	SortedSetKeys []string
	Tables        map[string]*Table
	// Logs          []ChangeLog
}

// FindTableByName finds a table by its name in the database
func (d *Database) FindTableByName(tableName string) (table *Table, e error) {
	// Search for table
	for _, dbTable := range d.Tables {
		if dbTable.Name == tableName {
			table = dbTable
			break
		}
	}

	if table == nil {
		e = errors.New("table not found")
	}

	return
}

// Table represents a table in a database
type Table struct {
	Name          string             `json:"name"`
	Engine        string             `json:"engine"`
	Version       int                `json:"version"`
	RowFormat     string             `json:"rowFormat"`
	Rows          int64              `json:"rows"`
	DataLength    int64              `json:"dataLength"`
	Collation     string             `json:"collation"`
	AutoIncrement int64              `json:"autoIncrement"`
	Columns       map[string]*Column `json:"columns"`
}

// Column represents a column in a table
type Column struct {
	Name       string `json:"column"`
	Position   int    `json:"position"`
	Default    string `json:"default"`
	IsNullable bool   `json:"isNullable"`
	DataType   string `json:"dataType"`
	MaxLength  int    `json:"maxLength"`
	Precision  int    `json:"precision"`
	CharSet    string `json:"charSet"`
	Type       string `json:"type"`
	ColumnKey  string `json:"columnKey"`
	Extra      string `json:"extra"`
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
