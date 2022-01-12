package schema

import (
	"database/sql"
	"errors"
	"sort"
	"strings"
)

const (
	RenameTable        = "RENAME_TABLE"
	CreateTable        = "CREATE_TABLE"
	DropTable          = "DROP_TABLE"
	ChangeColumn       = "CHANGE_COLUMN"
	AddColumn          = "ADD_COLUMN"
	DropColumn         = "DROP_COLUMN"
	AddIndex           = "ADD_INDEX"
	DropIndex          = "DROP_INDEX"
	ChangeCharacterSet = "CHANGE_CHARACTER_SET"
)

type SchemaList struct {
	Schemas []*Schema `json:"schemas"`
}

// Database represents a database
type Database struct {
	RunID               int64                               `json:"-"`
	Name                string                              `json:"name"`
	SortedSetKeys       []string                            `json:"-"`
	Tables              map[string]*Table                   `json:"tables"`
	Enums               map[string][]map[string]interface{} `json:"-"`
	DefaultCharacterSet string                              `json:"defaultCharacterSet"`
	DefaultCollation    string                              `json:"defaultCollation"`
}

func (d *Database) ToSchema(schemaName string) *Schema {
	return &Schema{
		Name:                schemaName,
		SortedSetKeys:       d.SortedSetKeys,
		Tables:              d.Tables,
		Enums:               d.Enums,
		DefaultCharacterSet: d.DefaultCharacterSet,
		DefaultCollation:    d.DefaultCollation,
	}
}

// Schema represents a database structure
type Schema struct {
	RunID               int64                               `json:"-"`
	Name                string                              `json:"name"`
	SortedSetKeys       []string                            `json:"-"`
	Tables              map[string]*Table                   `json:"tables"`
	Enums               map[string][]map[string]interface{} `json:"-"`
	DefaultCharacterSet string                              `json:"defaultCharacterSet"`
	DefaultCollation    string                              `json:"defaultCollation"`
}

// ToSortedTables returns SortedTables
func (s *Schema) ToSortedTables() SortedTables {

	sortedTables := make(SortedTables, 0, len(s.Tables))

	for _, table := range s.Tables {
		sortedTables = append(sortedTables, table)
	}

	sort.Sort(sortedTables)

	return sortedTables
}

// FindTableByName finds a table by its name in the database
func (s *Schema) FindTableByName(tableName string) (table *Table, e error) {
	// Search for table
	for _, dbTable := range s.Tables {
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
	CharacterSet  string             `json:"characterSet"`
	AutoIncrement int64              `json:"-"`
	Columns       map[string]*Column `json:"columns"`
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
	return strings.Compare(c[i].Name, c[j].Name) < 0
	// return c[i].Position < c[j].Position
}

// ToSortedColumns returns SortedColumns
func (table *Table) ToSortedColumns() SortedColumns {

	sortedColumns := make(SortedColumns, 0, len(table.Columns))

	// Find the primary key
	for _, column := range table.Columns {
		sortedColumns = append(sortedColumns, column)
	}

	sort.Sort(sortedColumns)

	return sortedColumns
}

// SortedTables is a slice of Table objects
type SortedTables []*Table

// Len is part of sort.Interface.
func (c SortedTables) Len() int {
	return len(c)
}

// Swap is part of sort.Interface.
func (c SortedTables) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// Less is part of sort.Interface. We use count as the value to sort by
func (c SortedTables) Less(i, j int) bool {
	// return c[i].Name < c[j].Name
	return strings.Compare(c[i].Name, c[j].Name) < 0
}

// Column represents a column in a table
type Column struct {
	Name string `json:"column"`
	// Position     int    `json:"position"`
	Default      string `json:"default"`
	IsNullable   bool   `json:"isNullable"`
	IsUnsigned   bool   `json:"isUnsigned"`
	DataType     string `json:"dataType"`
	MaxLength    int    `json:"maxLength"`
	Precision    int    `json:"precision"`
	CharSet      string `json:"charSet"`
	Collation    string `json:"collation"`
	Type         string `json:"type"`
	ColumnKey    string `json:"columnKey"`
	NumericScale int    `json:"numericScale"`
	Extra        string `json:"extra"`
	FmtType      string `json:"fmtType"`
	GoType       string `json:"goType"`
	IsString     bool   `json:"isString"`
}

// ColumnWithTable is a column with the table name included
type ColumnWithTable struct {
	*Column
	TableName string
}

// Server represents a server
type Server struct {
	Name            string `json:"name"`
	Host            string `json:"host"`
	Databases       map[string]*Schema
	Connection      *sql.DB
	CurrentDatabase string
}

// SchemaExists checks if the database `databaseName` exists in its list of databases
func (s *Server) SchemaExists(schemaName string) bool {

	exists := false

	for _, db := range s.Databases {
		if db.Name == schemaName {
			exists = true
			break
		}
	}

	return exists
}
