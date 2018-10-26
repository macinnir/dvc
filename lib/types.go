package lib

import (
	"database/sql"
	"errors"
)

// Options are the available runtime flags
type Options uint

// Command is the command line functionality
type Command struct {
	Options Options
}

// Changeset represents all of the changes in an environment and their changes
type Changeset struct {
	ChangeFiles map[string]ChangeFile
	Signature   string
}

// ChangeFile represents both a physical file on the local file system
// along with the entry in the changefile database
type ChangeFile struct {
	ID          int64
	DateCreated int64
	Hash        string
	Name        string
	DVCSetID    int64
	IsRun       bool
	IsDeleted   bool
	Content     string
	FullPath    string
	Ordinal     int
}

// Config contains a set of configuration values used throughout the application
type Config struct {
	Host          string `toml:"host"`
	DatabaseName  string `toml:"databaseName"`
	Username      string `toml:"username"`
	Password      string `toml:"password"`
	ChangeSetPath string `toml:"changesetPath"`
	DatabaseType  string `toml:"databaseType"`
	ReposDir      string `toml:"reposDir"`
	CachesDir     string `toml:"cachesDir"`
	ModelsDir     string `toml:"modelsDir"`
	TypescriptDir string `toml:"typescriptDir"`
	SchemaDir     string `toml:"schemaDir"`
}

// SortedColumns is a slice of Column objects
type SortedColumns []*Column

// Len is part of sort.Interface.
func (c SortedColumns) Len() int {
	return len(c)
}

// Swap is part of sort.Interface.
func (c SortedColumns) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// Less is part of sort.Interface. We use count as the value to sort by
func (c SortedColumns) Less(i, j int) bool {
	return c[i].Position < c[j].Position
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
	Host  string `json:"-"`
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
	Rows          int64              `json:"-"`
	DataLength    int64              `json:"-"`
	Collation     string             `json:"collation"`
	AutoIncrement int64              `json:"-"`
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

const (
	// OptLogInfo triggers verbose logging
	OptLogInfo = 1 << iota
	// OptLogDebug triggers extremely verbose logging
	OptLogDebug
	// OptSilent suppresses all logging
	OptSilent
	// OptReverse reverses the function
	OptReverse
	// OptSummary shows a summary of the actions instead of a raw stdout dump
	OptSummary
	// OptClean cleans
	OptClean
)

// IConnector defines the shape of a connector to a database
type IConnector interface {
	ConnectToServer(host string, username string, password string) (server *Server, e error)
	FetchDatabases(server *Server) (databases map[string]*Database, e error)
	UseDatabase(server *Server, databaseName string) (e error)
	FetchDatabaseTables(server *Server, databaseName string) (tables map[string]*Table, e error)
	FetchTableColumns(server *Server, databaseName string, tableName string) (columns map[string]*Column, e error)
}
