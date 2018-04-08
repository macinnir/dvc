package main

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
}

// Database represents a database
type Database struct {
	RunID int64
	Name  string
	Host  string
	// Sets          map[string]*ChangeSet
	SortedSetKeys []string
	Tables        map[string]Table
	// Logs          []ChangeLog
}

// Table represents a table in a database
type Table struct {
	Name          string            `json:"name"`
	Engine        string            `json:"engine"`
	Version       int               `json:"version"`
	RowFormat     string            `json:"rowFormat"`
	Rows          int64             `json:"rows"`
	DataLength    int64             `json:"dataLength"`
	Collation     string            `json:"collation"`
	AutoIncrement int64             `json:"autoIncrement"`
	Columns       map[string]Column `json:"columns"`
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

// Server represents a server
type Server struct {
	Name      string `json:"name"`
	Host      string `json:"host"`
	Databases map[string]*Database
}

// Log reprents a log entry
type Log struct {
	LogType    string
	LogMessage string
}
